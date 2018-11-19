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
)

type Peer struct {
	PublicKey string
	Address   string
	BootNodes *sync.Map
	PeerList  *sync.Map
}

func (peer *Peer) Initialize(publicKey string, address string) {
	peer.PublicKey = publicKey
	peer.Address = address
	peer.BootNodes = new(sync.Map)
	peer.PeerList = new(sync.Map)

	// FIXME create a configuration file
	peer.BootNodes.Store("BN0", "127.0.0.1:3300")
	peer.BootNodes.Store("BN1", "127.0.0.1:3301")
}

//Run test
func (peer *Peer) Run() {
	go peer.getPeer()
	go peer.ping()
	peer.pong()
}

func (peer *Peer) getPeer() {
	for {
		peer.BootNodes.Range(func(publicKey, address interface{}) bool {
			time.Sleep(2 * time.Second)
			addressStr, _ := address.(string)
			conn, err := net.Dial("tcp", addressStr)

			if err == nil {
				jsonStr := marshal(peer)
				conn.Write([]byte(jsonStr))
				// log.Printf("Send: %s", jsonStr)

				buff := make([]byte, 1024)
				n, _ := conn.Read(buff)
				// log.Printf("Receive: %s", buff[:n])

				data := string(buff[:n])
				p := unmarshal([]byte(data))
				// log.Printf("Receive: %s", data)
				peer.updatePeerList(p.PeerList)
			}
			return true
		})
	}
}

// ping pings to peers in peerList periodically
func (peer *Peer) ping() {
	for {
		peer.PeerList.Range(func(publicKey, address interface{}) bool {
			time.Sleep(2 * time.Second)
			addressStr, _ := address.(string)
			conn, err := net.Dial("tcp", addressStr)

			if err == nil {
				str := addressStr
				conn.Write([]byte(peer.Address))
				log.Printf("Ping: %s", str)

				buff := make([]byte, 1024)
				n, _ := conn.Read(buff)
				log.Printf("Pong: %s", buff[:n])
			}
			return true
		})
	}
}

// pong listens and pongs to peers when receives ping
func (peer *Peer) pong() {
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
					log.Println("Received from: ", data)
					break RECEIVING

				default:
					log.Fatalf("Failed to receive data: %s", err)
					return
				}
			}

			writer.Write([]byte(peer.Address))
			writer.Flush()
			// log.Printf("Sent: %s", peer.Address)
		}(conn, peer)
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
