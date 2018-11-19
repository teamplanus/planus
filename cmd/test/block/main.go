package main

import (
	"github.com/teamplanus/planus/core"
)

func main() {
	blockchain := core.NewBlockchain()
	blockchain.AddBlock(blockchain.CreateBlock("New Block 1"))
	blockchain.AddBlock(blockchain.CreateBlock("New Block 2"))

	// Test
	blockchain.ShowBlockchainForDebug()
}
