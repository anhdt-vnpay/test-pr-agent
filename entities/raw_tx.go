package entities

import (
	"encoding/json"
	"time"
)

type RawTransaction struct {
	Payload                []byte    `json:"data"`
	Txhash                 string    `json:"tx_hash"`
	BlockTime              time.Time `json:"block_time"`
	ChaincodeName          string    `json:"cc_name"`
	ValidationCode         string    `json:"code"`
	ChaincodeProposalInput string    `json:"func"`
	BlockNum               uint64    `json:"bn"`
	NetworkName            string    `json:"nnn"`
	ChannelName            string    `json:"cnn"`
	LastUpdate             int64     `json:"ludp"`
}

func (t *RawTransaction) MarshalBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *RawTransaction) UnmarshalBytes(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *RawTransaction) GetItemKey() string {
	return t.Txhash
}

func (t *RawTransaction) SetLastUpdate(lastUpdate int64) error {
	t.LastUpdate = lastUpdate
	return nil
}

func (t *RawTransaction) GetLastUpdate() int64 {
	return t.LastUpdate
}

func (t *RawTransaction) IsNotNil() bool {
	return t != nil
}
