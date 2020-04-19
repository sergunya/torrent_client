package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/marksamman/bencode"
)

type TorrentClient struct {
	PeerID      [20]byte
	TorrentFile TorrentFile
	URL         string
}

func (t *TorrentClient) generatePeerID() error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])

	if err == nil {
		t.PeerID = peerID
	}

	return err
}

func (t *TorrentClient) buildURL() error {

	base, err := url.Parse(t.TorrentFile.Announce)
	if err != nil {
		return err
	}
	params := url.Values{
		"info_hash":  []string{string(t.TorrentFile.InfoHash[:])},
		"peer_id":    []string{string(t.PeerID[:])},
		"port":       []string{"6889"},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{string(t.TorrentFile.Length)},
	}
	base.RawQuery = params.Encode()
	t.URL = base.String()

	return nil
}

func (t *TorrentClient) getPeersInfo() ([]Peer, int64, error) {
	resp, err := http.Get(t.URL)
	if err != nil {
		log.Fatal("cannot build url")
	}

	decodedBody, err := bencode.Decode(resp.Body)
	if err != nil {
		log.Fatal("cannot decode response body")
	}

	peers, err := DecodePeerInfo([]byte(decodedBody["peers"].(string)))

	if err != nil {
		log.Fatal("cannot decode peers info")
	}

	return peers, decodedBody["interval"].(int64), nil
}

func (t *TorrentClient) connect(p Peer) net.Conn {
	conn, err := net.Dial("tcp", net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port))))
	if err != nil {
		log.Fatal("cannot connect to peer")
	}

	// Start handshake procedure
	_, err = conn.Write(SerializeHandshake(t.TorrentFile.InfoHash, t.PeerID))
	if err != nil {
		log.Fatal("cannot send handshake to my bro peer")
	}

	answerInfoHash, answerPeerID, err := readHandshakeAnswer(conn)
	if err != nil {
		log.Fatal("cannot read answer to handshake")
	}

	if answerInfoHash != t.TorrentFile.InfoHash {
		conn.Close()
		log.Fatal("answered infoHash is not equal with original hash")
	}

	return conn
}

func CreateClient(t TorrentFile) (TorrentClient, error) {

	torrentClient := TorrentClient{TorrentFile: t}
	err := torrentClient.generatePeerID()
	if err != nil {
		fmt.Println(err)
	}

	err = torrentClient.buildURL()
	if err != nil {
		fmt.Println(err)
	}

	return torrentClient, nil
}
