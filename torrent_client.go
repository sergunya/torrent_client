package main

import "crypto/rand"

type TorrentClient struct {
}

func generatePeerID() ([20]byte, error) {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])

	return peerID, err
}

func makeRequest(t TorrentFile) {

}
