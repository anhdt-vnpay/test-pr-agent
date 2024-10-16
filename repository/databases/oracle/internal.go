package oracle

import (
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/types"
)

type transform interface {
	BlockToBlockOracle(blocks []*entities.Block) *types.BlocksOracle
	RawTxsToRawTxsOracle(rawTxs []*entities.RawTransaction) *types.RawTransactionsOracle
	TasksToTaskOracle(tasks []*entities.Task) *types.TaskOracle
	OnchainTxsToOnchainTxOracle(transactions []*entities.OnchainTransaction) *types.OnchainTransactionsOracle
	AccountTxsToAccountTxOracle(accountTxs []*entities.AccountTx) *types.AccountTxsOracle
	AccountsToAccountOracle(accounts []*entities.Account) *types.AccountsOracle
}
