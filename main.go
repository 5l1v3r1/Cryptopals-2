package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	first()
	second()
}

func second() {
	a_hex := "1c0111001f010100061a024b53535009181c"
	b_hex := "686974207468652062756c6c277320657965"

	a, err := hex.DecodeString(a_hex)
	if err != nil {
		fmt.Printf("Some error: %v", err)
		return
	}

	b, err := hex.DecodeString(b_hex)
	if err != nil {
		fmt.Printf("Some error: %v", err)
		return
	}

	n := len(a)
	c := make([]byte, n)
	for i := 0; i < n; i++ {
		c[i] = a[i] ^ b[i]
	}
	c_hex := hex.EncodeToString(c)
	fmt.Printf("a: %s\nb: %s\nc: %s\n", a_hex, b_hex, c_hex)
}

func first() {
	hex_string := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	bytes, err := hex.DecodeString(hex_string)
	if err != nil {
		fmt.Printf("Some error: %v", err)
		return
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	fmt.Printf("Hex String: %s\n", hex_string)
	fmt.Printf("Base64 String: %s\n", b64)
}
