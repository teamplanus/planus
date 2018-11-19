package main

import (
	"os"

	"github.com/teamplanus/planus/p2p/client"
)

func main() {

	peer := new(p2p.Peer)
	peer.Initialize(os.Args[1], os.Args[2])
	peer.Run()
}
