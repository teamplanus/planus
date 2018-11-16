package core

import (
	"log"
)

// Blockchain represents a list of connected blocks
type Blockchain struct {
	 blocks []*Block
}

// NewBlockchain creates a blockchain with genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// AddBlock adds the block into blockchain received from other peers
func (blockchain *Blockchain) AddBlock(block *Block) {
	if (blockchain.isValidBlock(block)) {
		blockchain.blocks = append(blockchain.blocks, block)
	}
}

// CreateBlock creates a new block
func (blockchain *Blockchain) CreateBlock(data string) *Block {
	return NewBlock(
		blockchain.GetLatestBlockNumber() + 1,
		blockchain.GetLatestBlock().PreviousHash, data)
}

// GetLatestBlock returns latest block in the blockchain
func (blockchain *Blockchain) GetLatestBlock() *Block {
	return blockchain.blocks[blockchain.GetLatestBlockNumber()]
}

// GetLatestBlockNumber returns latest block number in the blockchain
func (blockchain *Blockchain) GetLatestBlockNumber() uint64 {
	// FIXME Use other methods instead of len() that supports unit64
	return (uint64)(len(blockchain.blocks) - 1)
}

// isValidBlock validates block and returns valid or not
func (blockchain *Blockchain) isValidBlock(block *Block) bool {
	// TODO Validate block number with GetLatestBlockNumber()
	// TODO Validate previous block hash with GetLatestBlock()
	// TODO Validate new block hash
	return true
}

// FOR TEST
func (blockchain *Blockchain) ShowBlockchainForDebug() {
	for idx, block := range blockchain.blocks {
		log.Println(idx, block.Data)
	}
}