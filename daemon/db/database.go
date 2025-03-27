package db

import (
	"embed"
	"fmt"
	"sync"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

//go:embed schema/*.sql
var embedSchema embed.FS

type SqliteDatabaseConfig struct {
	File         string `yaml:"file"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
}

type Database struct {
	config      *SqliteDatabaseConfig
	logger      logrus.FieldLogger
	ReaderDb    *sqlx.DB
	writerDb    *sqlx.DB
	writerMutex sync.Mutex
}

func NewDatabase(config *SqliteDatabaseConfig, logger logrus.FieldLogger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Init() error {
	if d.config.MaxOpenConns == 0 {
		d.config.MaxOpenConns = 50
	}
	if d.config.MaxIdleConns == 0 {
		d.config.MaxIdleConns = 10
	}
	if d.config.MaxOpenConns < d.config.MaxIdleConns {
		d.config.MaxIdleConns = d.config.MaxOpenConns
	}

	d.logger.Infof("initializing sqlite connection to %v with %v/%v conn limit", d.config.File, d.config.MaxIdleConns, d.config.MaxOpenConns)
	dbConn, err := sqlx.Open("sqlite", fmt.Sprintf("%s?_pragma=journal_mode(WAL)", d.config.File))
	if err != nil {
		return fmt.Errorf("error opening sqlite database: %v", err)
	}

	d.checkDbConn(dbConn, "database")
	dbConn.SetConnMaxIdleTime(0)
	dbConn.SetConnMaxLifetime(0)
	dbConn.SetMaxOpenConns(d.config.MaxOpenConns)
	dbConn.SetMaxIdleConns(d.config.MaxIdleConns)

	dbConn.MustExec("PRAGMA journal_mode = WAL")

	d.ReaderDb = dbConn
	d.writerDb = dbConn

	return nil
}

func (d *Database) Close() error {
	err := d.writerDb.Close()
	if err != nil {
		return fmt.Errorf("error closing writer db connection: %v", err)
	}
	return nil
}

func (d *Database) checkDbConn(dbConn *sqlx.DB, dataBaseName string) {
	// The golang sql driver does not properly implement PingContext
	// therefore we use a timer to catch db connection timeouts
	dbConnectionTimeout := time.NewTimer(15 * time.Second)

	go func() {
		<-dbConnectionTimeout.C
		d.logger.Fatalf("timeout while connecting to %s", dataBaseName)
	}()

	err := dbConn.Ping()
	if err != nil {
		d.logger.Fatalf("unable to Ping %s: %s", dataBaseName, err)
	}

	dbConnectionTimeout.Stop()
}

func (d *Database) RunDBTransaction(handler func(tx *sqlx.Tx) error) error {
	d.writerMutex.Lock()
	defer d.writerMutex.Unlock()

	tx, err := d.writerDb.Beginx()
	if err != nil {
		return fmt.Errorf("error starting db transactions: %v", err)
	}

	defer tx.Rollback()

	err = handler(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing db transaction: %v", err)
	}

	return nil
}

func (d *Database) ApplyEmbeddedDbSchema(version int64) error {
	goose.SetLogger(&gooseLogger{logger: d.logger})
	goose.SetBaseFS(embedSchema)
	schemaDirectory := "schema"

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if version == -2 {
		if err := goose.Up(d.writerDb.DB, schemaDirectory, goose.WithAllowMissing()); err != nil {
			return err
		}
	} else if version == -1 {
		if err := goose.UpByOne(d.writerDb.DB, schemaDirectory, goose.WithAllowMissing()); err != nil {
			return err
		}
	} else {
		if err := goose.UpTo(d.writerDb.DB, schemaDirectory, version, goose.WithAllowMissing()); err != nil {
			return err
		}
	}

	return nil
}

type gooseLogger struct {
	logger logrus.FieldLogger
}

func (g *gooseLogger) Fatalf(format string, v ...interface{}) {
	g.logger.Fatalf(format, v...)
}

func (g *gooseLogger) Printf(format string, v ...interface{}) {
	g.logger.Infof(format, v...)
}
