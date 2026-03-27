package maelstromx

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Decode[T any](msg maelstrom.Message) (T, error) {
	var body T
	err := json.Unmarshal(msg.Body, &body)
	return body, err
}
