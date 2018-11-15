package core

type Block struct {
	
	BlockNumber		uint64
	Hash			[]byte
	PreviousHash	[]byte
	TimeStamp		uint64	// Unix timestamp (https://www.unixtimestamp.com/)
	Data			string	// FIXME It will be changed to some struct contains txs.
}