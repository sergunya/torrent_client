package main

import "fmt"

func main() {
	path := "torrent_files/debian.torrent"
	tFile, _ := ProcessFile(path)
	tClient, _ := CreateClient(tFile)
	peers, interval, _ := tClient.getPeersInfo()
	tClient.connect(peers[0])
	fmt.Println(interval)
	//tracker.TorrentClient.connect(tracker.Peers[0])

}
