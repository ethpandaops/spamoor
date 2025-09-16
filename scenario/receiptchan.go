package scenario

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiptChan chan *types.Receipt

func (rc ReceiptChan) Wait(ctx context.Context) (*types.Receipt, error) {
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
