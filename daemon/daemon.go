package daemon

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

type Daemon struct {
	ctx        context.Context
	logger     logrus.FieldLogger
	clientPool *spamoor.ClientPool
	rootWallet *txbuilder.Wallet
	txpool     *txbuilder.TxPool
	db         *db.Database

	spammerIdMtx  sync.Mutex
	spammerMap    map[int64]*Spammer
	spammerMapMtx sync.RWMutex
}

func NewDaemon(ctx context.Context, logger logrus.FieldLogger, clientPool *spamoor.ClientPool, rootWallet *txbuilder.Wallet, txpool *txbuilder.TxPool, db *db.Database) *Daemon {
	return &Daemon{
		ctx:        ctx,
		logger:     logger,
		clientPool: clientPool,
		rootWallet: rootWallet,
		txpool:     txpool,
		db:         db,
		spammerMap: make(map[int64]*Spammer),
	}
}

func (d *Daemon) GetClientPool() *spamoor.ClientPool {
	return d.clientPool
}

func (d *Daemon) Run() error {
	// restore all spammer from db
	spammerList, err := d.db.GetSpammers()
	if err != nil {
		return fmt.Errorf("failed to get all spammer: %w", err)
	}

	for _, spammer := range spammerList {
		_, err := d.restoreSpammer(spammer)
		if err != nil {
			return fmt.Errorf("failed to restore spammer: %w", err)
		}
	}

	return nil
}

func (d *Daemon) GetSpammer(id int64) *Spammer {
	d.spammerMapMtx.RLock()
	defer d.spammerMapMtx.RUnlock()

	return d.spammerMap[id]
}

func (d *Daemon) GetAllSpammers() []*Spammer {
	d.spammerMapMtx.RLock()
	defer d.spammerMapMtx.RUnlock()

	spammerList := make([]*Spammer, 0, len(d.spammerMap))
	for _, spammer := range d.spammerMap {
		spammerList = append(spammerList, spammer)
	}

	sort.Slice(spammerList, func(i, j int) bool {
		return spammerList[i].GetID() > spammerList[j].GetID()
	})

	return spammerList
}

func (d *Daemon) DeleteSpammer(id int64) error {
	d.spammerMapMtx.Lock()
	defer d.spammerMapMtx.Unlock()

	spammer := d.spammerMap[id]
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}

	// Stop if running
	if spammer.scenarioCancel != nil {
		spammer.Pause()
	}

	// Delete from DB
	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.DeleteSpammer(tx, id)
	})
	if err != nil {
		return fmt.Errorf("failed to delete spammer: %w", err)
	}

	// Remove from map
	delete(d.spammerMap, id)
	return nil
}
