package transform

import (
	"github.com/blcvn/corev3-libs/ledger/hyperledger/helper"
	"github.com/blcvn/corev4-explorer/entities"
)

func (t *transform) TransformBlockStateProto(ret *helper.TxState, shardId string, blockNumber uint64, txIndex, txCount int) (*entities.RawTransaction, error) {
	tx := &entities.RawTransaction{
		Payload:                ret.Data,
		Txhash:                 ret.TxId,
		ChaincodeName:          ret.Chaincode,
		ValidationCode:         ret.ValidationCode.String(),
		ChaincodeProposalInput: ret.FuncName,
		BlockNum:               blockNumber,
		NetworkName:            ret.Channel,
		ChannelName:            ret.Channel,
	}
	return tx, nil
}
