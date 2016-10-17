package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
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

// Takes two byte slices and computes the hammind distance of their bits.
// If one slice is shorter than the other, the shorter will be extended with 0s.
func Hamming(buf1, buf2 []byte) int {
	score := 0

	// Make both buffers same size padded by 0s
	if len(buf1) > len(buf2) {
		dif := len(buf1) - len(buf2)
		new_buf := make([]byte, len(buf2)+dif)
		for i, b := range buf2 {
			new_buf[i] = b
		}
		buf2 = new_buf
	} else if len(buf1) < len(buf2) {
		dif := len(buf2) - len(buf1)
		new_buf := make([]byte, len(buf1)+dif)
		for i, b := range buf1 {
			new_buf[i] = b
		}
		buf1 = new_buf
	}

	// I'm glad I remember my bitwise math from college
	for i, b := range buf1 {
		x := b ^ buf2[i]
		d, _ := binary.ReadUvarint(bytes.NewReader([]byte{x}))
		if d >= 128 {
			score++
			d = d - 128
		}
		if d >= 64 {
			score++
			d = d - 64
		}
		if d >= 32 {
			score++
			d = d - 32
		}
		if d >= 16 {
			score++
			d = d - 16
		}
		if d >= 8 {
			score++
			d = d - 8
		}
		if d >= 4 {
			score++
			d = d - 4
		}
		if d >= 2 {
			score++
			d = d - 2
		}
		if d >= 1 {
			score++
			d = d - 1
		}
	}
	return score
}
