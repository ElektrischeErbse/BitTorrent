package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var readIndex int = 0

func decodeString(bencode string) (string, error) {
	// 字符串 <长度>:<字符串内容>
	// - 5:hello -> hello
	// - 10:hello12345 -> hello12345
	var firstColonIndex int = -1
	var startIndex = readIndex
	for ; readIndex < len(bencode); readIndex++ {
		if rune(bencode[readIndex]) == ':' {
			firstColonIndex = readIndex
			break
		}
	}
	if firstColonIndex == -1 {
		return "", fmt.Errorf("Parse string failed: syntax error, expected character ':'")
	}
	length, err := strconv.Atoi(bencode[startIndex:firstColonIndex])
	if err != nil {
		return "", fmt.Errorf("Parse string failed: %v", err)
	}
	readIndex = readIndex + length + 1
	return bencode[firstColonIndex+1 : firstColonIndex+1+length], nil
}

func decodeInteger(bencode string) (int, error) {
	// 整数 以 i 开头，e 结尾
	// - i42e -> 42
	var integerEnd int = -1
	var startIndex = readIndex + 1
	for ; readIndex < len(bencode); readIndex++ {
		if rune(bencode[readIndex]) == 'e' {
			integerEnd = readIndex
			break
		}
	}
	if integerEnd == -1 {
		return 0, fmt.Errorf("Parse integer failed: syntax error, expected character 'e'")
	}
	num, err := strconv.Atoi(bencode[startIndex:integerEnd])
	if err != nil {
		return 0, fmt.Errorf("Parse integer failed: %v", err)
	}
	readIndex = integerEnd + 1
	return num, nil
}

func decodeList(bencode string) ([]any, error) {
	// 列表 以 l 开头，e 结尾
	// - l4:spam4:eggse -> ["spam", "eggs"]
	// - l4:spami100ee -> ["spam", 100]
	// - l5:applel6:bananai42eee -> ["apple", ["banana", 42]]
	var items []any
	readIndex++
	for rune(bencode[readIndex]) != 'e' {
		item, err := DecodeBencode(bencode)
		if err != nil {
			return items, fmt.Errorf("Parse list failed: %v", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func decodeDict(bencode string) (map[string]any, error) {
	// 字典 以 d 开头，e 结尾，键值对交替出现
	// - d3:agei25e4:name5:alicee -> {"age": 25, "name": "alice"}
	// - d4:tagsl1:a1:bee -> {"tags": ["a", "b"]}
	// - d6:configd5:debugi1eee -> {"config": {"debug": 1}}
	var dict = make(map[string]any)
	readIndex++
	if !unicode.IsDigit(rune(bencode[readIndex])) {
		return dict, fmt.Errorf("Parse dict failed: key is not a string")
	}
	for rune(bencode[readIndex]) != 'e' {
		key, err := decodeString(bencode)
		if err != nil {
			return dict, fmt.Errorf("Parse dict failed: %v", err)
		}
		value, err := DecodeBencode(bencode)
		if err != nil {
			return dict, fmt.Errorf("Parse dict failed: %v", err)
		}
		dict[key] = value
	}
	return dict, nil
}

func DecodeBencode(bencode string) (any, error) {
	if unicode.IsDigit(rune(bencode[readIndex])) {
		return decodeString(bencode)
	} else if rune(bencode[readIndex]) == 'i' {
		return decodeInteger(bencode)
	} else if rune(bencode[readIndex]) == 'l' {
		return decodeList(bencode)
	} else if rune(bencode[readIndex]) == 'd' {
		return decodeDict(bencode)
	} else {
		return nil, fmt.Errorf("Invalid bencode!")
	}
}

func main() {
	command := os.Args[1]
	switch command {
	case "decode":
		bencodeValue := os.Args[2]
		result, err := DecodeBencode(bencodeValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		jsonOutput, _ := json.Marshal(result)
		fmt.Println(string(jsonOutput))
	case "info":
		torrentFile := os.Args[2]
		t, err := ParseTorrentFile(torrentFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse torrent file: %s failed, %v\n", torrentFile, err)
			os.Exit(1)
		}
		fmt.Printf("Tracker: %s\nLength: %d\n", t.Tracker, t.Length)
		fmt.Printf("Info Hash: %s\n", hex.EncodeToString(t.InfoHash[:]))
		fmt.Printf("Piece Length: %d\n", t.PieceLength)
		fmt.Printf("Piece Hashes:\n")
		for i := 0; i < len(t.PieceHashes); i += 20 {
			fmt.Println(hex.EncodeToString(t.PieceHashes[i : i+20]))
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
