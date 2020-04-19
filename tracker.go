package main

type Tracker struct {
	Interval      int64
	Peers         []Peer
	TorrentClient TorrentClient
}
