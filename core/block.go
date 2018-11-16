package core

// Block represents a block in the blockchain
type Block struct {
	BlockNumber		uint64
	Hash			[]byte
	PreviousHash	[]byte
	TimeStamp		uint64	// Unix timestamp (https://www.unixtimestamp.com/)
	Data			string	// FIXME It will be changed to some struct contains txs.
}

// CreateBlock creates and returns new block
func CreateBlock(blockNumber uint64, previousHash []byte, data string) *Block {
	// TODO
	return &Block{}
}

// CreateGenesisBlock creates and returns genesis block
func CreateGenesisBlock() *Block {
	return CreateBlock(0, []byte{}, "Planus")
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