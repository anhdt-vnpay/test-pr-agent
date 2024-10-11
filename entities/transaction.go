package entities

import "time"

type Transaction struct {
	OriginTraceNo         string //Revert TraceNo
	TraceNo               string
	SenderId              string
	ReceiverId            string
	Description           string
	SenderBalanceBefore   int64
	SenderBalanceAfter    int64
	SenderOverdraftBefore int64
	SenderOverdraftAfter  int64
	ReceiverBalanceBefore int64
	ReceiverBalanceAfter  int64
	CreatedTime           int64 // time in milliseconds, generate from offchain
	HoldUntil             int64
	FeeAccountId          string
	SenderAmount          int64
	ReceiverAmount        int64
	FeeAmount             int64
	Txnhash               string
	TransformTaskId       string
}

type RawTransaction struct {
	Payload                []byte
	Txhash                 string
	BlockTime              time.Time
	ChaincodeName          string
	ValidationCode         string
	ChaincodeProposalInput string
	BlockNum               uint64
	NetworkName            string
	ChannelName            string
	FunctionName           string
	TxData                 string
}

type OnchainTransaction struct {
	Txhash          string
	TraceNo         string
	OriginTraceNo   string
	SenderId        string
	ReceiverId      string
	FeeAccountId    string
	SenderAmount    int64
	ReceiverAmount  int64
	FeeAmount       int64
	BlockNum        uint64
	NetworkName     string
	ChannelName     string
	TxType          string
	TimeOnchain     time.Time
	TimeOffchain    time.Time
	PreDay          time.Time
	CurrentDay      time.Time
	NextDay         time.Time
	EndOfMonth      string
	TransformTaskId int64
}
