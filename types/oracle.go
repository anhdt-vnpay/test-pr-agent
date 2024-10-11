package types

import "time"

type AccountsOracle struct {
	TraceNo         []string
	Txhash          []string
	AccountId       []string
	Amount          []int64
	BlockadeAmount  []int64
	Balance         []int64
	Type            []int32
	State           []int32
	TimeOnchain     []time.Time
	TimeOffchain    []time.Time
	NetworkName     []string
	ChannelName     []string
	BlockNum        []uint64
	TransformTaskId []int64
}

type AccountTxsOracle struct {
	TraceNo      []string
	Txhash       []string
	AccountId    []string
	Amount       []int64
	TimeOnchain  []time.Time
	TimeOffchain []time.Time
	CreatedAt    []time.Time
	NetworkName  []string
	ChannelName  []string
	BlockNum     []uint64

	PreDay          []time.Time
	NextDay         []time.Time
	CurrentDay      []time.Time
	EndOfMonth      []string
	TransformTaskId []string
}

type AccountsBalanceOracle struct {
	AccountId   []string
	Balance     []int64
	ProcessTime []time.Time
	NetworkName []string
	ChannelName []string
	TaskId      []string
}

type AccountsDeltaOracle struct {
	AccountId      []string
	Amount         []int64
	ProcessHour    []time.Time
	TimeCurrentDay []time.Time
	DeltaHour      []int
	DeltaDate      []time.Time
	TaskId         []string
}

type BlocksOracle struct {
	Blocknum           []uint64
	Datahash           []string
	Prehash            []string
	Txcount            []int64
	BlockTime          []time.Time
	Blockhash          []string
	PrevBlockhash      []string
	ChannelGenesisHash []string
	Blksize            []int64
	NetworkName        []string
	ChannelName        []string
}

type TaskOracle struct {
	Id          []int64
	Type        []int32
	Status      []int32
	BlockNumber []uint64
	NetworkName []string
	ChannelName []string
	CreateAt    []time.Time
}

type OnchainTransactionsOracle struct {
	Txhash          []string
	TraceNo         []string
	OriginTraceNo   []string
	SenderId        []string
	ReceiverId      []string
	FeeAccountId    []string
	SenderAmount    []int64
	ReceiverAmount  []int64
	FeeAmount       []int64
	BlockNum        []uint64
	NetworkName     []string
	ChannelName     []string
	TxType          []string
	TimeOnchain     []time.Time
	TimeOffchain    []time.Time
	PreDay          []time.Time
	NextDay         []time.Time
	CurrentDay      []time.Time
	EndOfMonth      []string
	TransformTaskId []int64
}

type RawTransactionsOracle struct {
	// Payload                []byte
	Txhash                 []string
	BlockTime              []time.Time
	ChaincodeName          []string
	ValidationCode         []string
	ChaincodeProposalInput []string
	BlockNum               []uint64
	NetworkName            []string
	ChannelName            []string
}
