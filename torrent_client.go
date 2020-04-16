package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

func (t *TorrentClient) makeRequest() (Tracker, error) {
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
		return Tracker{}, err
	}

	tracker := Tracker{
		Interval: decodedBody["interval"].(int64),
		Peers:    peers,
	}

	return tracker, nil
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
