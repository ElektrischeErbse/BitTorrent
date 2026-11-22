package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func decodeBencode(bencode string) (any, error) {
	if unicode.IsDigit(rune(bencode[0])) {
		// 字符串
		// - 5:hello -> hello
		// - 10:hello12345 -> hello12345
		var firstColonIndex int
		for i, c := range bencode {
			if c == ':' {
				firstColonIndex = i
				break
			}
		}
		length, err := strconv.Atoi(bencode[:firstColonIndex])
		if err != nil {
			return "", err
		}
		return bencode[firstColonIndex+1 : firstColonIndex+1+length], nil
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
