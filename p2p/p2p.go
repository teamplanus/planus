package p2p

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Peer struct {
	PublicKey string
	Address   string
	BootNodes *sync.Map
	PeerList  *sync.Map
}

func (peer *Peer) Initialize(publicKey string, address string, bootPubKey string, bootAddr string) {
	peer.PublicKey = publicKey
	peer.Address = address
	peer.BootNodes = new(sync.Map)
	peer.PeerList = new(sync.Map)

	// FIXME create a configuration file
	MASTER_NODE_PUBLICKEY := bootPubKey
	MASTER_NODE_ADDRESS := bootAddr

	peer.BootNodes.Store(MASTER_NODE_PUBLICKEY, MASTER_NODE_ADDRESS)
}

func (peer *Peer) Run() {
	go peer.ping()
	go peer.pong()
	peer.getPeers()
}

// ping pings to peers in peerList periodically
func (peer *Peer) ping() {
	for {
		peer.BootNodes.Range(func(publicKey, address interface{}) bool {
			publicKeyStr, _ := publicKey.(string)
			if publicKeyStr != peer.PublicKey {
				time.Sleep(2 * time.Second)
				addressStr, _ := address.(string)
				conn, err := net.Dial("tcp", addressStr+"0")

				if err == nil {
					jsonStr := marshal(peer)
					conn.Write([]byte(jsonStr))

					color.Cyan("Send: %s\n", jsonStr)
				}
			}
			return true
		})
	}
}

// pong listens and pongs to peers when receives ping
func (peer *Peer) pong() {
	listen, err := net.Listen("tcp4", peer.Address+"0")
	defer listen.Close()

	if err != nil {
		log.Fatalf("Failed to listen %s, %s", peer.Address+"0", err)
		os.Exit(1)
	}
	log.Printf("Start to listen %s", peer.Address+"0")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go func(conn net.Conn, peer *Peer) {
			defer conn.Close()

			var (
				buf    = make([]byte, 1024)
				reader = bufio.NewReader(conn)
				writer = bufio.NewWriter(conn)
			)
		RECEIVING:
			for {
				n, err := reader.Read(buf)
				data := string(buf[:n])

				switch err {
				case io.EOF:
					break RECEIVING
				case nil:
					color.Blue("Receive: %s\n", data)
					p := unmarshal([]byte(data))
					peer.updatePeerList(p.PeerList)

					break RECEIVING

				default:
					log.Fatalf("Failed to receive data: %s", err)
					return
				}
			}
			jsonStr := marshal(peer)
			writer.Write([]byte(jsonStr))
			writer.Flush()
		}(conn, peer)
	}
}

func (peer *Peer) getPeers() {
	listen, err := net.Listen("tcp4", peer.Address)
	defer listen.Close()

	if err != nil {
		log.Fatalf("Failed to listen %s, %s", peer.Address, err)
		os.Exit(1)
	}
	log.Printf("Start to listen %s", peer.Address)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go func(conn net.Conn, peer *Peer) {
			defer conn.Close()

			var (
				buf    = make([]byte, 1024)
				reader = bufio.NewReader(conn)
				writer = bufio.NewWriter(conn)
			)
		RECEIVING:
			for {
				n, err := reader.Read(buf)
				data := string(buf[:n])

				switch err {
				case io.EOF:
					break RECEIVING
				case nil:
					log.Println("Receive: ", data)

					p := unmarshal([]byte(data))
					peer.updatePeer(p)

					break RECEIVING
				default:
					log.Fatalf("Failed to receive data: %s", err)
					return
				}
			}
			jsonStr := marshal(peer)
			writer.Write([]byte(jsonStr))
			writer.Flush()
		}(conn, peer)
	}
}

func (peer *Peer) updatePeer(newPeer Peer) {
	if _, exists := peer.PeerList.Load(newPeer.PublicKey); !exists {
		log.Println("Update new peer: " + newPeer.PublicKey + " / " + newPeer.Address)
		peer.PeerList.Store(newPeer.PublicKey, newPeer.Address)
	}
}

func (peer *Peer) updatePeerList(newPeerList *sync.Map) {
	newPeerList.Range(func(publicKey, address interface{}) bool {
		publicKeyStr, _ := publicKey.(string)
		addressStr, _ := address.(string)

		if _, exists := peer.PeerList.Load(publicKeyStr); !exists {
			log.Println("Update new peer: " + publicKeyStr + " / " + addressStr)
			peer.PeerList.Store(publicKey, addressStr)
		}
		return true
	})
}

func marshal(peer *Peer) string {
	marshaled, _ := json.Marshal(peer)
	return strings.Replace(string(marshaled), "\"PeerList\":{}", string(marshalSyncMap(peer.PeerList, "PeerList")), -1)
}

func marshalSyncMap(syncMap *sync.Map, keyName string) []byte {
	buffer := bytes.NewBufferString("\"" + keyName + "\":{")
	count := 0
	length := 0
	syncMap.Range(func(_, _ interface{}) bool {
		length++
		return true
	})

	syncMap.Range(func(key, value interface{}) bool {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return false
		}
		buffer.WriteString(fmt.Sprintf("\"%s\":%s", key, string(jsonValue)))
		count++
		if count < length {
			buffer.WriteString(",")
		}
		return true
	})
	buffer.WriteString("}")
	return buffer.Bytes()
}

func unmarshal(data []byte) Peer {
	type _Peer struct {
		PublicKey string
		Address   string
		PeerList  map[string]string
	}

	var _p _Peer
	_ = json.Unmarshal(data, &_p)

	var p Peer
	p.PublicKey = _p.PublicKey
	p.Address = _p.Address
	p.BootNodes = new(sync.Map)
	p.PeerList = new(sync.Map)
	for key, value := range _p.PeerList {
		p.PeerList.Store(key, value)
	}
	return p
}
