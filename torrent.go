package main

import (
	"fmt"
	"os"
)

type Torrent struct {
	Tracker string
	Length  int
}

func ParseTorrentFile(torrentFile string) (Torrent, error) {
	var torrent Torrent
	bs, err := os.ReadFile(torrentFile)
	if err != nil {
		return torrent, err
	}
	content := string(bs)
	result, err := decodeBencode(content)
	if err != nil {
		return torrent, err
	}
	tracker, ok := getValueByKey[string](result, "announce")
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to get announce")
	}
	torrent.Tracker = tracker
	length, ok := getValueByKey[int](result, "info", "length")
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to get info.length")
	}
	torrent.Length = length
	return torrent, nil
}

func getValueByKey[T any](root any, keys ...string) (T, bool) {
	var zero T
	current := root
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return zero, false
		}
		val, exist := m[key]
		if !exist {
			return zero, false
		}
		current = val
	}
	result, ok := current.(T)
	return result, ok
}
