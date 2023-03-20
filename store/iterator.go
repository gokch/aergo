package store

import (
	"github.com/aergoio/aergo-lib/db"
	"github.com/aergoio/aergo/types"
)

type BlockTx struct {
	Block *types.Block
	Txs   []*types.Transaction
}

func IterateTransactions(db db.DB, from, to uint64, interrupt chan struct{}) chan *BlockTx {
	if from >= to {
		return nil
	}
	return nil
}
