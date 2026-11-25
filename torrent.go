package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
)

type Torrent struct {
	Tracker     string
	Length      int
	InfoHash    string
	PieceLength int
	PieceHashes []string
}

func ParseTorrentFile(torrentFile string) (Torrent, error) {
	var torrent Torrent
	bs, err := os.ReadFile(torrentFile)
	if err != nil {
		return torrent, err
	}
	content := string(bs)
	result, err := DecodeBencode(content)
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
		fmt.Fprintln(os.Stderr, "Failed to get info[length]")
	}
	torrent.Length = length
	info, ok := getValueByKey[map[string]any](result, "info")
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to get info")
	}
	infoContent, err := BuildBencode(info)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuildBencode error: %v\n", err)
	}
	hash := sha1.Sum([]byte(infoContent))
	torrent.InfoHash = hex.EncodeToString(hash[:])
	pieceLen, ok := getValueByKey[int](result, "info", "piece length")
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to get info[piece length]")
	}
	torrent.PieceLength = pieceLen
	// 所有数据块的 SHA-1 哈希值拼接而成的二进制串
	pieceHashes, ok := getValueByKey[string](result, "info", "pieces")
	if !ok {
		fmt.Fprintln(os.Stderr, "Failed to get info[pieces]")
	}
	for i := 0; i < len(pieceHashes); i += 20 {
		hexHash := hex.EncodeToString([]byte(pieceHashes[i : i+20]))
		torrent.PieceHashes = append(torrent.PieceHashes, hexHash)
	}
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

func BuildBencode(data any) (string, error) {
	switch v := data.(type) {
	case string:
		return encodeString(v)
	case int:
		return encodeInteger(v)
	case []any:
		return encodeList(v)
	case map[string]any:
		return encodeDict(v)
	default:
		return "", fmt.Errorf("Error: unsupported type!")
	}
}

func encodeString(data string) (string, error) {
	return fmt.Sprintf("%d:%s", len(data), data), nil
}

func encodeInteger(data int) (string, error) {
	return fmt.Sprintf("i%de", data), nil
}

func encodeList(data []any) (string, error) {
	res := "l"
	for _, item := range data {
		s, err := BuildBencode(item)
		if err != nil {
			return res, err
		}
		res += s
	}
	res += "e"
	return res, nil
}

func encodeDict(data map[string]any) (string, error) {
	// 键按字典序升序排列
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := "d"
	for _, k := range keys {
		s, err := encodeString(k)
		if err != nil {
			return res, err
		}
		res += s
		s, err = BuildBencode(data[k])
		if err != nil {
			return res, err
		}
		res += s
	}
	res += "e"
	return res, nil
}
