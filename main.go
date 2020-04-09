package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var PEER_ID = []byte{0xf0, 0x6f, 0x0f, 0xdf, 0x59, 0x7a, 0x12, 0x85, 0x63, 0x36, 0x95, 0xf5, 0x0a, 0x77, 0x3c, 0x91, 0xd9, 0xf1, 0xa4, 0xe5}

var NODE_ID = []byte{0xac, 0x79, 0x59, 0xbe, 0xdf, 0x58, 0xb7, 0x88, 0x2e, 0x61, 0x08, 0x10, 0x39, 0x1d, 0x2d, 0xf9, 0x80, 0xcb, 0xbb, 0xea}

func main() {
	log.Print("Main start")

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	// 1. Start tcp peer
	peerAddress, err := net.ResolveTCPAddr("tcp", "localhost:3333")
	if err != nil {
		log.Panic(err)
	}
	peer, err := net.ListenTCP("tcp", peerAddress)
	if err != nil {
		log.Panic(err)
	}

	// 2. Start udp node
	nodeAddress, err := net.ResolveUDPAddr("udp", "localhost:3334")
	if err != nil {
		log.Panic(err)
	}

	node, err := net.ListenUDP("udp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}

	// 3. Run peer and node
	go run_peer(peer)
	go run_node(node)

	// 4. Setup shutdown
	go graceful_shutdown(peer, node, quit, done)

	<-done
	log.Println("Main exit")
}

func graceful_shutdown(peer *net.TCPListener, node *net.UDPConn, quit chan os.Signal, done chan bool) {
	<-quit
	log.Println("Gracefully shutting down...")
	peer.Close()
	node.Close()

	close(done)
}

const MAX_CONNECTIONS = 32

// BitTorrent TCP peer
func run_peer(peer *net.TCPListener) {
	conn, err := peer.AcceptTCP()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	go func(conn net.Conn) {

	}(conn)
}

const NODE_TYPE_GOOD = 0x1
const NODE_TYPE_QUESTIONABLE = 0x2

//
// Torrent or metainfo file
// For the Torrent file spec @see http://www.bittorrent.org/beps/bep_0003.html#metainfo-files 09.04.2020
//
type Torrent struct {
	// The URL of the tracker.
	announce string

	// This maps to a dictionary, with keys described below.
	info Info
}

// Torrent info dictionary
type Info struct {
	// The name key maps to a UTF-8 encoded string
	// which is the suggested name to save the file (or directory) as.
	// It is purely advisory.
	name string

	// most commonly 2^18 = 256K
	piece_length int

	// Pieces maps to a string whose length is a multiple of 20.
	// It is to be subdivided into strings of length 20,
	// each of which is the SHA1 hash of the piece at the corresponding index.
	pieces []string

	// There is also a key length or a key files, but not both or neither.
	// If length is present then the download represents a single file,
	// otherwise it represents a set of files which go in a directory structure.
	//
	// In the single file case, length maps to the length of the file in bytes.
	length int

	// For the purposes of the other keys, the multi-file case is treated as only having a
	// single file by concatenating the files in the order they appear in the files list.
	// The files list is the value files maps to, and is a list of dictionaries
	// containing the following keys:
	files File
}

// Torrent file in a torrent directory
type File struct {
	// length - The length of the file, in bytes.
	length int

	// path - A list of UTF-8 encoded strings corresponding to subdirectory names,
	// the last of which is the actual file name
	// (a zero length list is an error case).
	path string
}

type Node struct {
	node_type    int
	last_changed time.Time
}

var routing_table []Node

// DHT UDP node
func run_node(node *net.UDPConn) {
	var buf [1024]byte

	// Reader
	go func() {
		for {
			_, _, err := node.ReadFromUDP(buf[:])
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}
