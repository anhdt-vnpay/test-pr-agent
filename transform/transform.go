package transform

import (
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/types"
)

type transform struct{}

func NewTransform() *transform {
	return &transform{}
}

func (t *transform) AccountsToAccountOracle(accounts []*entities.Account) *types.AccountsOracle {
	accountOracle := &types.AccountsOracle{}
	for _, account := range accounts {
		accountOracle.TraceNo = append(accountOracle.TraceNo, account.TraceNo)
		accountOracle.Txhash = append(accountOracle.Txhash, account.Txhash)
		accountOracle.AccountId = append(accountOracle.AccountId, account.AccountId)
		accountOracle.BlockadeAmount = append(accountOracle.BlockadeAmount, account.BlockadeAmount)
		accountOracle.Balance = append(accountOracle.Balance, account.Balance)
		accountOracle.Type = append(accountOracle.Type, account.Type)
		accountOracle.State = append(accountOracle.State, account.State)
		accountOracle.TimeOnchain = append(accountOracle.TimeOnchain, account.TimeOnchain)
		accountOracle.TimeOffchain = append(accountOracle.TimeOffchain, account.TimeOffchain)
		accountOracle.NetworkName = append(accountOracle.NetworkName, account.NetworkName)
		accountOracle.ChannelName = append(accountOracle.ChannelName, account.ChannelName)
		accountOracle.BlockNum = append(accountOracle.BlockNum, account.BlockNum)
		accountOracle.TransformTaskId = append(accountOracle.TransformTaskId, account.TransformTaskId)
	}
	return accountOracle
}

func (t *transform) AccountTxsToAccountTxOracle(accountTxs []*entities.AccountTx) *types.AccountTxsOracle {
	accountTxsOracle := &types.AccountTxsOracle{}
	for _, accountTx := range accountTxs {
		accountTxsOracle.TraceNo = append(accountTxsOracle.TraceNo, accountTx.TraceNo)
		accountTxsOracle.Txhash = append(accountTxsOracle.Txhash, accountTx.Txhash)
		accountTxsOracle.AccountId = append(accountTxsOracle.AccountId, accountTx.AccountId)
		accountTxsOracle.Amount = append(accountTxsOracle.Amount, accountTx.Amount)
		accountTxsOracle.TimeOnchain = append(accountTxsOracle.TimeOnchain, accountTx.TimeOnchain)
		accountTxsOracle.TimeOffchain = append(accountTxsOracle.TimeOffchain, accountTx.TimeOffchain)
		accountTxsOracle.CreatedAt = append(accountTxsOracle.CreatedAt, accountTx.CreatedAt)
		accountTxsOracle.NetworkName = append(accountTxsOracle.NetworkName, accountTx.NetworkName)
		accountTxsOracle.ChannelName = append(accountTxsOracle.ChannelName, accountTx.ChannelName)
		accountTxsOracle.BlockNum = append(accountTxsOracle.BlockNum, accountTx.BlockNum)
		accountTxsOracle.PreDay = append(accountTxsOracle.PreDay, accountTx.PreDay)
		accountTxsOracle.NextDay = append(accountTxsOracle.NextDay, accountTx.NextDay)
		accountTxsOracle.CurrentDay = append(accountTxsOracle.CurrentDay, accountTx.CurrentDay)
		accountTxsOracle.EndOfMonth = append(accountTxsOracle.EndOfMonth, accountTx.EndOfMonth)
		accountTxsOracle.TransformTaskId = append(accountTxsOracle.TransformTaskId, accountTx.TransformTaskId)
	}
	return accountTxsOracle
}

func (t *transform) AccountsBalanceToAccountBalanceOracle(accountBalances []*entities.AccountsBalance) *types.AccountsBalanceOracle {
	accountsBalanceOracle := &types.AccountsBalanceOracle{}
	for _, accountBalance := range accountBalances {
		accountsBalanceOracle.AccountId = append(accountsBalanceOracle.AccountId, accountBalance.AccountId)
		accountsBalanceOracle.Balance = append(accountsBalanceOracle.Balance, accountBalance.Balance)
		accountsBalanceOracle.ProcessTime = append(accountsBalanceOracle.ProcessTime, accountBalance.ProcessTime)
		accountsBalanceOracle.NetworkName = append(accountsBalanceOracle.NetworkName, accountBalance.NetworkName)
		accountsBalanceOracle.ChannelName = append(accountsBalanceOracle.ChannelName, accountBalance.ChannelName)
		accountsBalanceOracle.TaskId = append(accountsBalanceOracle.TaskId, accountBalance.TaskId)
	}
	return accountsBalanceOracle
}

func (t *transform) AccountDeltaToAccountDeltaOracle(accountsDelta []*entities.AccountsDelta) *types.AccountsDeltaOracle {
	accountsDeltaOracle := &types.AccountsDeltaOracle{}
	for _, accountDelta := range accountsDelta {
		accountsDeltaOracle.AccountId = append(accountsDeltaOracle.AccountId, accountDelta.AccountId)
		accountsDeltaOracle.Amount = append(accountsDeltaOracle.Amount, accountDelta.Amount)
		accountsDeltaOracle.ProcessHour = append(accountsDeltaOracle.ProcessHour, accountDelta.ProcessHour)
		accountsDeltaOracle.TimeCurrentDay = append(accountsDeltaOracle.TimeCurrentDay, accountDelta.TimeCurrentDay)
		accountsDeltaOracle.DeltaHour = append(accountsDeltaOracle.DeltaHour, accountDelta.DeltaHour)
		accountsDeltaOracle.DeltaDate = append(accountsDeltaOracle.DeltaDate, accountDelta.DeltaDate)
		accountsDeltaOracle.TaskId = append(accountsDeltaOracle.TaskId, accountDelta.TaskId)
	}
	return accountsDeltaOracle
}

func (t *transform) BlockToBlockOracle(blocks []*entities.Block) *types.BlocksOracle {
	blockOracle := &types.BlocksOracle{}
	for _, block := range blocks {
		blockOracle.Blocknum = append(blockOracle.Blocknum, block.Blocknum)
		blockOracle.Datahash = append(blockOracle.Datahash, block.Datahash)
		blockOracle.Prehash = append(blockOracle.Prehash, block.Prehash)
		blockOracle.Txcount = append(blockOracle.Txcount, block.Txcount)
		blockOracle.BlockTime = append(blockOracle.BlockTime, block.BlockTime)
		blockOracle.Blockhash = append(blockOracle.Blockhash, block.Blockhash)
		blockOracle.PrevBlockhash = append(blockOracle.PrevBlockhash, block.PrevBlockhash)
		blockOracle.ChannelGenesisHash = append(blockOracle.ChannelGenesisHash, block.ChannelGenesisHash)
		blockOracle.Blksize = append(blockOracle.Blksize, block.Blksize)
		blockOracle.NetworkName = append(blockOracle.NetworkName, block.NetworkName)
		blockOracle.ChannelName = append(blockOracle.ChannelName, block.ChannelName)
	}
	return blockOracle
}

func (t *transform) OnchainTxsToOnchainTxOracle(transactions []*entities.OnchainTransaction) *types.OnchainTransactionsOracle {
	ochainTxsOracle := &types.OnchainTransactionsOracle{}
	for _, tx := range transactions {
		ochainTxsOracle.Txhash = append(ochainTxsOracle.Txhash, tx.Txhash)
		ochainTxsOracle.TraceNo = append(ochainTxsOracle.TraceNo, tx.TraceNo)
		ochainTxsOracle.OriginTraceNo = append(ochainTxsOracle.OriginTraceNo, tx.OriginTraceNo)
		ochainTxsOracle.SenderId = append(ochainTxsOracle.SenderId, tx.SenderId)
		ochainTxsOracle.ReceiverId = append(ochainTxsOracle.ReceiverId, tx.ReceiverId)
		ochainTxsOracle.FeeAccountId = append(ochainTxsOracle.FeeAccountId, tx.FeeAccountId)
		ochainTxsOracle.SenderAmount = append(ochainTxsOracle.SenderAmount, tx.SenderAmount)
		ochainTxsOracle.ReceiverAmount = append(ochainTxsOracle.ReceiverAmount, tx.ReceiverAmount)
		ochainTxsOracle.FeeAmount = append(ochainTxsOracle.FeeAmount, tx.FeeAmount)
		ochainTxsOracle.BlockNum = append(ochainTxsOracle.BlockNum, tx.BlockNum)
		ochainTxsOracle.NetworkName = append(ochainTxsOracle.NetworkName, tx.NetworkName)
		ochainTxsOracle.ChannelName = append(ochainTxsOracle.ChannelName, tx.ChannelName)
		ochainTxsOracle.TxType = append(ochainTxsOracle.TxType, tx.TxType)
		ochainTxsOracle.TimeOnchain = append(ochainTxsOracle.TimeOnchain, tx.TimeOnchain)
		ochainTxsOracle.TimeOffchain = append(ochainTxsOracle.TimeOffchain, tx.TimeOffchain)
		ochainTxsOracle.PreDay = append(ochainTxsOracle.PreDay, tx.PreDay)
		ochainTxsOracle.CurrentDay = append(ochainTxsOracle.CurrentDay, tx.CurrentDay)
		ochainTxsOracle.NextDay = append(ochainTxsOracle.NextDay, tx.NextDay)
		ochainTxsOracle.EndOfMonth = append(ochainTxsOracle.EndOfMonth, tx.EndOfMonth)
		ochainTxsOracle.TransformTaskId = append(ochainTxsOracle.TransformTaskId, tx.TransformTaskId)
	}
	return ochainTxsOracle
}

func (t *transform) RawTxsToRawTxsOracle(rawTxs []*entities.RawTransaction) *types.RawTransactionsOracle {
	rawTxOracle := &types.RawTransactionsOracle{}
	for _, rawTx := range rawTxs {
		rawTxOracle.Txhash = append(rawTxOracle.Txhash, rawTx.Txhash)
		rawTxOracle.BlockTime = append(rawTxOracle.BlockTime, rawTx.BlockTime)
		rawTxOracle.ChaincodeName = append(rawTxOracle.ChaincodeName, rawTx.ChaincodeName)
		rawTxOracle.ValidationCode = append(rawTxOracle.ValidationCode, rawTx.ValidationCode)
		rawTxOracle.ChaincodeProposalInput = append(rawTxOracle.ChaincodeProposalInput, rawTx.ChaincodeProposalInput)
		rawTxOracle.NetworkName = append(rawTxOracle.NetworkName, rawTx.NetworkName)
		rawTxOracle.ChannelName = append(rawTxOracle.ChannelName, rawTx.ChannelName)
		rawTxOracle.BlockNum = append(rawTxOracle.BlockNum, rawTx.BlockNum)
	}
	return rawTxOracle
}

func (t *transform) TasksToTaskOracle(tasks []*entities.Task) *types.TaskOracle {
	tasksOracle := &types.TaskOracle{}
	for _, task := range tasks {
		tasksOracle.Id = append(tasksOracle.Id, task.Id)
		tasksOracle.Type = append(tasksOracle.Type, task.Type)
		tasksOracle.Status = append(tasksOracle.Status, task.Status)
		tasksOracle.BlockNumber = append(tasksOracle.BlockNumber, task.BlockNumber)
		tasksOracle.NetworkName = append(tasksOracle.NetworkName, task.NetworkName)
		tasksOracle.ChannelName = append(tasksOracle.ChannelName, task.ChannelName)

	}
	return tasksOracle
}
