package core

// Block represents a block in the blockchain
type Block struct {
	BlockNumber		uint64
	Hash			[]byte
	PreviousHash	[]byte
	TimeStamp		uint64	// Unix timestamp (https://www.unixtimestamp.com/)
	Data			string	// FIXME It will be changed to some struct contains txs.
}

// NewBlock creates and returns new block
func NewBlock(blockNumber uint64, previousHash []byte, data string) *Block {
	// TODO
	return &Block{blockNumber, []byte{}, previousHash, 1234567890, data}
}

// NewGenesisBlock creates and returns genesis block
func NewGenesisBlock() *Block {
	return NewBlock(0, []byte{}, "Planus Genesis Block")
}

// Serialize serializes block to byte array
func (block *Block) Serialize() []byte {
	// TODO
	return []byte{}
}

// Deserialize deserializes byte array to block
func Deserialize(bytes []byte) *Block {
	// TODO
	return &Block{}
}