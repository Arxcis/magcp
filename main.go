package main

import (
	"log"
	"net"
	"os"
	"os/signal"
)

var PEER_ID = []byte{0xf0, 0x6f, 0x0f, 0xdf, 0x59, 0x7a, 0x12, 0x85, 0x63, 0x36, 0x95, 0xf5, 0x0a, 0x77, 0x3c, 0x91, 0xd9, 0xf1, 0xa4, 0xe5}

var NODE_ID = []byte{0xac, 0x79, 0x59, 0xbe, 0xdf, 0x58, 0xb7, 0x88, 0x2e, 0x61, 0x08, 0x10, 0x39, 0x1d, 0x2d, 0xf9, 0x80, 0xcb, 0xbb, 0xea}

func main() {
	log.Print("Main start")

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	peerAddress, err := net.ResolveTCPAddr("tcp", "localhost:3333")
	if err != nil {
		log.Panic(err)
	}
	peer, err := net.ListenTCP("tcp", peerAddress)
	if err != nil {
		log.Panic(err)
	}

	nodeAddress, err := net.ResolveUDPAddr("udp", "localhost:3334")
	if err != nil {
		log.Panic(err)
	}

	node, err := net.ListenUDP("udp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}

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

// BitTorrent TCP peer
func run_peer(peer *net.Listener) {

}

// DHT UDP node
func run_node(node *net.Listener) {

}
