package main

import "fmt"

func main() {
	path := "torrent_files/debian.torrent"
	t, err := ProcessFile(path)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(t)
	}
}
