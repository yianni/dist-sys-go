package txn

import (
	"encoding/json"
	"fmt"
)

type operation struct {
	Kind  string
	Key   int
	Value *int
}

func (o *operation) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 3 {
		return fmt.Errorf("operation length = %d, want 3", len(raw))
	}

	if err := json.Unmarshal(raw[0], &o.Kind); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &o.Key); err != nil {
		return err
	}
	if string(raw[2]) == "null" {
		return nil
	}

	var value int
	if err := json.Unmarshal(raw[2], &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

func (o operation) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{o.Kind, o.Key, o.Value})
}

type txnRequest struct {
	Type string      `json:"type"`
	Txn  []operation `json:"txn"`
}

type txnResponse struct {
	Type string      `json:"type"`
	Txn  []operation `json:"txn"`
}

type syncRequest struct {
	Type   string       `json:"type"`
	Writes []writeState `json:"writes"`
}

type syncResponse struct {
	Type string `json:"type"`
}
