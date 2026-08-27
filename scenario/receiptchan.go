package scenario

import (
	"context"

	"github.com/ethpandaops/spamoor/txtypes"
)

type ReceiptChan chan *txtypes.Receipt

func (rc ReceiptChan) Wait(ctx context.Context) (*txtypes.Receipt, error) {
	if rc == nil {
		return nil, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case receipt := <-rc:
		return receipt, nil
	}
}
