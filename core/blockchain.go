package core

// Blockchain represents a list of connected blocks
type Blockchain struct {
	 blocks []*Block
}

// NewBlockchain creates a blockchain with genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// AddBlock adds the block into the blockchain
func (blockchain *Blockchain) AddBlock(block *Block) {
	// TODO
}