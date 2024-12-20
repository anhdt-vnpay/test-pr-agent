package entities

import "time"

type Block struct {
	Blocknum           uint64
	Datahash           string
	Prehash            string
	Txcount            int64
	BlockTime          time.Time
	Blockhash          string
	PrevBlockhash      string
	ChannelGenesisHash string
	Blksize            int64
	NetworkName        string
	ChannelName        string
}
