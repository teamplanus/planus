package main

import (
	"os"

	"github.com/teamplanus/planus/p2p/bootnode"
)

func main() {

	peer := new(p2p.Peer)
	peer.Initialize(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	peer.Run()
}
