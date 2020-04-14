package main

import (
	"fmt"
	"log"
)

func main() {
	path := "torrent_files/debian.torrent"
	t, err := ProcessFile(path)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(t)
		client, err := CreateClient(t)
		if err != nil {
			log.Fatal(err)
		}

		err = client.makeRequest()
		fmt.Println("Stop")
	}
}
