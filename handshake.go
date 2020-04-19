package main

import (
	"io"
)

const preposition = "BitTorrent protocol"

func SerializeHandshake(infoHash [20]byte, peerID [20]byte) []byte {
	buf := make([]byte, len(preposition)+49)
	buf[0] = byte(len(preposition))
	curr := 1
	curr += copy(buf[curr:], preposition)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], infoHash[:])
	curr += copy(buf[curr:], peerID[:])
	return buf
}

func readHandshakeAnswer(r io.Reader) ([20]byte, [20]byte, error) {
	var infoHash, peerID [20]byte

	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)

	if err != nil {
		return infoHash, peerID, err
	}

	prepLen := int(lengthBuf[0])

	handshakeBuf := make([]byte, 48+prepLen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return infoHash, peerID, err
	}

	copy(infoHash[:], handshakeBuf[prepLen+8:prepLen+8+20])
	copy(peerID[:], handshakeBuf[prepLen+8+20:])

	return infoHash, peerID, nil
}
