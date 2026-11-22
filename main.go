package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var readIndex int = 0

func decodeBencode(bencode string) (any, error) {
	if unicode.IsDigit(rune(bencode[readIndex])) {
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
	} else if rune(bencode[readIndex]) == 'i' {
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
	} else if rune(bencode[readIndex]) == 'l' {
		// 列表 以 l 开头，e 结尾
		// - l4:spam4:eggse -> ["spam", "eggs"]
		// - l4:spami100ee -> ["spam", 100]
		// - l5:applel6:bananai42eee -> ["apple", ["banana", 42]]
		var items []any
		readIndex++
		for rune(bencode[readIndex]) != 'e' {
			item, err := decodeBencode(bencode)
			if err != nil {
				return items, fmt.Errorf("Parse list failed: %v", err)
			}
			items = append(items, item)
		}
		return items, nil
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func main() {
	command := os.Args[1]
	if command == "decode" {
		bencodeValue := os.Args[2]
		result, err := decodeBencode(bencodeValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		jsonOutput, _ := json.Marshal(result)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
