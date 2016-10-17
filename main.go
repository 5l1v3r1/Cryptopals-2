package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func HexToBase64(hex_string string) (string, error) {
	buf, err := hex.DecodeString(hex_string)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func XOR(buf1, buf2 []byte) ([]byte, error) {
	if len(buf1) != len(buf2) {
		return nil, errors.New("XOR: slices must be equal size")
	}
	out := make([]byte, len(buf1))
	for i, b := range buf1 {
		out[i] = b ^ buf2[i]
	}
	return out, nil
}

// Buf will be XOR'd with key with repeating key if necessary.
func RepeatingXOR(buf, key []byte) []byte {
	if len(key) == 0 {
		return buf
	}
	div := len(buf) / len(key)
	rem := len(buf) % len(key)

	rkey := make([]byte, len(buf))
	for i := 0; i < div; i++ {
		for j := 0; j < len(key); j++ {
			rkey[len(key)*i+j] = key[j]
		}
	}
	for i := 0; i < rem; i++ {
		rkey[len(key)*div+i] = key[i]
	}

	// Won't have error as both inputs will be equal length
	res, _ := XOR(buf, rkey)
	return res
}

func SingleCharOracle(ct []byte) string {
	keyspace := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	results := make(map[string]string)

	for i := 0; i < len(keyspace); i++ {
		k := keyspace[i]
		results[string(k)] = string(RepeatingXOR(ct, []byte{k}))
	}

	highscore := 0
	highkey := ""
	for k, res := range results {
		score := PtScore(res)
		if score > highscore {
			highscore = score
			highkey = k
		}
	}
	return highkey
}

// Plaintext scoring algorithm
func PtScore(pt string) int {
	score := 0
	for _, c := range pt {
		// A-Z
		if c >= 65 && c <= 90 {
			score++
		}
		// a-z
		if c >= 97 && c <= 122 {
			score++
		}
		// non-printable
		if c < 32 || c == 127 {
			score--
		}
	}
	return score
}
