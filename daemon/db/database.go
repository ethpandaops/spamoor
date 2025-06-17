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

// SqliteDatabaseConfig defines the configuration for SQLite database connections.
// It specifies the database file path and connection pool limits for managing
// concurrent database access efficiently.
type SqliteDatabaseConfig struct {
	File         string `yaml:"file"`         // Database file path
	MaxOpenConns int    `yaml:"maxOpenConns"` // Maximum number of open connections to the database
	MaxIdleConns int    `yaml:"maxIdleConns"` // Maximum number of idle connections in the pool
}

// Database manages SQLite database connections with WAL mode and connection pooling.
// It provides both reader and writer connections with mutex protection for write operations
// and embedded schema migration support using goose.
type Database struct {
	config      *SqliteDatabaseConfig
	logger      logrus.FieldLogger
	ReaderDb    *sqlx.DB   // Database connection for read operations
	writerDb    *sqlx.DB   // Database connection for write operations
	writerMutex sync.Mutex // Protects write transactions from concurrent access
}

// NewDatabase creates a new Database instance with the specified configuration and logger.
// The database connections are not initialized until Init() is called.
func NewDatabase(config *SqliteDatabaseConfig, logger logrus.FieldLogger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

// Init initializes the database connections with WAL mode and connection pooling.
// Sets default connection limits (50 max open, 10 max idle) if not specified.
// Enables WAL mode for better concurrent access and configures connection timeouts.
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

// Close closes the database writer connection.
// Should be called during application shutdown to ensure proper cleanup.
func (d *Database) Close() error {
	err := d.writerDb.Close()
	if err != nil {
		return fmt.Errorf("error closing writer db connection: %v", err)
	}
	return nil
}

// checkDbConn validates database connectivity with a 15-second timeout.
// Uses a timer-based approach to catch connection timeouts since the SQL driver
// doesn't properly implement PingContext. Fatally exits if connection fails.
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

// RunDBTransaction executes a function within a database transaction with automatic rollback.
// The transaction is protected by a mutex to ensure sequential write operations.
// Automatically rolls back on error and commits on success.
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

// ApplyEmbeddedDbSchema applies database schema migrations using embedded SQL files.
// Supports different migration strategies: -2 (all), -1 (one up), or specific version.
// Uses goose migration library with allowMissing option for flexible schema management.
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

// gooseLogger adapts logrus.FieldLogger to the goose logger interface.
// Provides logging integration for database migration operations.
type gooseLogger struct {
	logger logrus.FieldLogger
}

// Fatalf logs a fatal message and exits the application.
// Implements the goose logger interface for fatal errors during migrations.
func (g *gooseLogger) Fatalf(format string, v ...interface{}) {
	g.logger.Fatalf(format, v...)
}

// Printf logs an informational message.
// Implements the goose logger interface for migration progress logging.
func (g *gooseLogger) Printf(format string, v ...interface{}) {
	g.logger.Infof(format, v...)
}
