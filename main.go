package main

import (
	"log"
	"./core"
)

func main() {
	blockchain := core.NewBlockchain()
	log.Println(blockchain.GetLastBlockNumber())
}