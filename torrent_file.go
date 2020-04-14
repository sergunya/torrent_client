package main

import (
	"crypto/sha1"
	"log"
	"os"

	"github.com/marksamman/bencode"
)

// Structure of parsed torrent file
// Almost all fields is obvious (InfoHash is sha1 of encoded `info` part of torrent file)
type TorrentFile struct {
	LocalFilePath string
	Announce      string
	InfoHash      [20]byte
	Length        int64
	Name          string
}

// Method of TorrentFile struct
// The main purpose is parse torrent file and fill all fields
func (t *TorrentFile) parseTorrentFile() error {

	file, err := os.Open(t.LocalFilePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer file.Close()

	result, err := bencode.Decode(file)
	if err != nil {
		log.Fatal(err)
		return err
	}

	t.Announce = result["announce"].(string)
	t.InfoHash = sha1.Sum(bencode.Encode(result["info"]))
	t.Length = result["info"].(map[string]interface{})["length"].(int64)
	t.Name = result["info"].(map[string]interface{})["name"].(string)

	return nil
}

func ProcessFile(path string) (TorrentFile, error) {
	if _, err := os.Stat(path); err == nil {

		myTorrentFile := TorrentFile{LocalFilePath: path}
		if err := myTorrentFile.parseTorrentFile(); err != nil {
			return TorrentFile{}, err
		}

		return myTorrentFile, nil

	} else {
		return TorrentFile{}, err
	}
}
