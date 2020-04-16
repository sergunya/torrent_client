package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func DecodePeerInfo(encodedPeers []byte) ([]Peer, error) {
	// 4 bytes for IP and 2 bytes for port
	peerSize := 6
	numPeers := len(encodedPeers) / peerSize

	if len(encodedPeers)%peerSize != 0 {
		log.Fatal("peers binary info is invalid")
		return nil, fmt.Errorf("peers binary info is invalid")
	}

	decodedPeers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		indent := i * peerSize
		p := Peer{
			IP:   net.IP(encodedPeers[indent : indent+4]),
			Port: binary.BigEndian.Uint16(encodedPeers[indent+4 : indent+6]),
		}

		decodedPeers[i] = p
	}

	return decodedPeers, nil
}
