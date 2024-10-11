package hyperledger

import "encoding/json"

type Tx struct {
	TxId       string `json:"tx_id"`
	Channel    string `json:"channel"`
	Chaincode  string `json:"chaincode"`
	FuncName   string `json:"func"`
	Data       []byte `json:"data"`
	LastUpdate int64
}

func (t *Tx) MarshalBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Tx) UnmarshalBytes(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *Tx) GetItemKey() string {
	return t.TxId
}

func (t *Tx) SetLastUpdate(lastUpdate int64) error {
	t.LastUpdate = lastUpdate
	return nil
}

func (t *Tx) GetLastUpdate() int64 {
	return t.LastUpdate
}

func (t *Tx) IsNotNil() bool {
	return t != nil
}
