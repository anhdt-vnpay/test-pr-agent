package hyperledger

import (
	"github.com/blcvn/corev3-libs/ledger/hyperledger/helper"
)

type transform struct{}

func (t *transform) TransformBlockStateProto(ret *helper.TxState, shardId string, blockNumber uint64, txIndex, txCount int) (*Tx, error) {
	tx := &Tx{
		TxId:      ret.TxId,
		Channel:   ret.Channel,
		Chaincode: ret.Chaincode,
		FuncName:  ret.FuncName,
		Data:      ret.Data,
	}
	return tx, nil
}
