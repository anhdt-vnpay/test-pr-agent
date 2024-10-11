package entities

import "time"

type Account struct {
	Txhash          string
	TraceNo         string
	AccountId       string
	Amount          int64
	BlockadeAmount  int64
	Balance         int64
	Type            int32
	State           int32
	TimeOnchain     time.Time
	TimeOffchain    time.Time
	NetworkName     string
	ChannelName     string
	BlockNum        uint64
	TransformTaskId string
}

type AccountTx struct {
	Txhash          string
	TraceNo         string
	AccountId       string
	Amount          int64
	TimeOnchain     time.Time
	TimeOffchain    time.Time
	NetworkName     string
	ChannelName     string
	BlockNum        uint64
	PreDay          time.Time
	NextDay         time.Time
	CurrentDay      time.Time
	EndOfMonth      string
	CreatedAt       time.Time
	TransformTaskId string
}

type AccountsBalance struct {
	AccountId   string
	Balance     int64
	ProcessTime time.Time
	NetworkName string
	ChannelName string
	TaskId      string
}

type AccountsDelta struct {
	AccountId      string
	Amount         int64
	ProcessHour    time.Time
	TimeCurrentDay time.Time
	DeltaHour      int
	DeltaDate      time.Time
	TaskId         string
}
