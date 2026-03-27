package maelstromx

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func MustRun(node *maelstrom.Node) {
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
