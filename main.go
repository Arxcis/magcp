package main

import (
	"log"
	"net"
	"os"
	"os/signal"
)

var PEER_ID = []byte{0xf0, 0x6f, 0x0f, 0xdf, 0x59, 0x7a, 0x12, 0x85, 0x63, 0x36, 0x95, 0xf5, 0x0a, 0x77, 0x3c, 0x91, 0xd9, 0xf1, 0xa4, 0xe5}

type TrackerRequest struct {
	// The 20 byte sha1 hash of the bencoded form of the info value from the metainfo file.
	// This value will almost certainly have to be escaped.
	info_hash string

	// A string of length 20 which this downloader uses as its id.
	// Each downloader generates its own id at random at the start of a new download.
	// This value will also almost certainly have to be escaped.
	peer_id string

	// An optional parameter giving the IP (or dns name) which this peer is at.
	// Generally used for the origin if it's on the same machine as the tracker.
	ip string

	// The port number this peer is listening on.
	// Common behavior is for a downloader to try to listen on port 6881
	// and if that port is taken try 6882, then 6883, etc. and give up after 6889.
	port string

	// The total amount uploaded so far, encoded in base ten ascii.
	uploaded int

	// The total amount downloaded so far, encoded in base ten ascii.
	downloaded int

	// The number of bytes this peer still has to download, encoded in base ten ascii.
	// Note that this can't be computed from downloaded
	// and the file length since it might be a resume,
	// and there's a chance that some of the downloaded data
	// failed an integrity check and had to be re-downloaded.
	left int

	// This is an optional key which maps to started, completed, or stopped
	// (or empty, which is the same as not being present).
	// If not present, this is one of the announcements done at regular intervals.
	// An announcement using started is sent when a download first begins,
	// and one using completed is sent when the download is complete.
	// No completed is sent if the file was complete when started.
	// Downloaders send an announcement using stopped when they cease downloading.
	event string
}

func main() {
	log.Print()
	log.Print("Main start")

	if len(os.Args) != 2 {
		log.Fatal("len(os.Args) != 2")
	}

	arg := os.Args[1]
	if arg != "magnet" && arg != "seed" {
		log.Fatal(`arg != "magnet" || arg != "seed"`)
	}

	if arg == "magnet" {
		log.Print("You are a magnet")
	} else {
		log.Print("You are a seeder")
		_, err := net.Dial("tcp", "localhost:3334")
		if err != nil {
			log.Fatal(err)
		}

	}

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	// 1. Start tcp peer
	peer, err := net.Listen("tcp", "localhost:3333")
	if err != nil {
		log.Fatal(err)
	}

	// 3. Run peer
	go run_peer(peer)

	// 4. Setup shutdown
	go graceful_shutdown(peer, quit, done)

	<-done
	log.Println("Main exit")
}

func graceful_shutdown(peer net.Listener, quit chan os.Signal, done chan bool) {
	<-quit
	log.Println("Gracefully shutting down...")
	peer.Close()
	close(done)
}

const MAX_CONNECTIONS = 32

// BitTorrent TCP peer
func run_peer(peer net.Listener) {
	for {
		conn, err := peer.Accept()
		if err != nil {
			log.Print(err)
			break
		}

		go run_connection(conn)
	}
}

func run_connection(conn net.Conn) {
	defer conn.Close()

}

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
