package main

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	keyspace := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+-=[]{};:'\",./<>?`~ "
	results := make(map[string]string)

	for i := 0; i < len(keyspace); i++ {
		k := keyspace[i]
		results[string(k)] = string(RepeatingXOR(ct, []byte{k}))
	}

	highscore := 9999.0
	highkey := ""
	for k, res := range results {
		score := PtScore(res)
		//fmt.Printf("String: %s\nScore: %f\n", strconv.Quote(string(RepeatingXOR(ct, []byte(k)))), score)
		if score < highscore {
			highscore = score
			highkey = k
		}
	}
	fmt.Printf("highkey: %s\n", highkey)
	return highkey
}

var english_freq = []float64{
	0.08167, 0.01492, 0.02782, 0.04253, 0.12702, 0.02228, 0.02015, // A-G
	0.06094, 0.06966, 0.00153, 0.00772, 0.04025, 0.02406, 0.06749, // H-N
	0.07507, 0.01929, 0.00095, 0.05987, 0.06327, 0.09056, 0.02758, // O-U
	0.00978, 0.02360, 0.00150, 0.01974, 0.00074, // V-Z
}

// Plaintext scoring algorithm
func PtScore(pt string) float64 {
	pt = strings.ToLower(pt)
	ptb := []byte(pt)
	count := make(map[byte]int)
	ignored := 0
	for _, c := range ptb {
		if c >= 97 && c <= 122 {
			count[c] = count[c] + 1
		} else {
			ignored++
		}
		if c >= 0 && c <= 8 {
			return 99999999
		}
		if c >= 11 && c <= 13 {
			return 99999999
		}
		if c >= 14 && c <= 31 {
			return 99999999
		}
		if c == 127 {
			return 99999999
		}
	}

	chi2 := 0.0
	for k, v := range count {
		expected := float64(len(pt)-ignored) * english_freq[k-97]
		dif := float64(v) - expected
		chi2 = chi2 + (dif * dif / expected)
	}

	return chi2
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

func KeySizeOracle(ct []byte) int {
	results := make(map[int]float64)

	for ks := 2; ks <= 40; ks++ {

		div := len(ct) / ks
		div = div - 2
		tot := 0.0
		for i := 0; i < div; i++ {
			a, b := ct[ks*i:ks*(i+1)], ct[ks*(i+1):ks*(i+2)]
			tot = tot + float64(Hamming(a, b))/float64(ks)
		}
		dif := tot / float64(div)

		//fmt.Printf("KS: %d, res: %f\n", ks, dif)
		results[ks] = dif

	}

	low := 999.9
	lowks := 0
	for k, v := range results {
		fmt.Printf("Keysize: %d, result: %f\n", k, v)
		if v < float64(low) {
			low = v
			lowks = k
		}
	}

	ks := lowks
	fmt.Printf("Keysize: %d\n", ks)

	return ks
}

func RepeatingKeyOracle(ct []byte) (pt []byte) {
	ks := KeySizeOracle(ct)

	blocks := make([][]byte, 0)
	div := len(ct) / ks
	fmt.Printf("div: %d\n", div)

	for i := 0; i < ks; i++ {
		var block []byte
		for j := 0; j < div; j++ {
			block = append(block, ct[j*ks+i])
		}
		fmt.Printf("Block %d: len %d\n", i, len(block))
		blocks = append(blocks, block)
	}
	fmt.Printf("Blocks %d\n", len(blocks))

	var key string
	for i, block := range blocks {
		c := SingleCharOracle(block)
		key = key + c
		if i == 24 {
			fmt.Printf("Block %d decrypted with %s: %s\n", i, c, strconv.Quote(string(RepeatingXOR(block, []byte(c)))))
		}
	}
	fmt.Printf("Key: %s", key)
	return RepeatingXOR(ct, []byte(key))
}

func ECBDecrypt(key, ct []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ct) < aes.BlockSize {
		panic("ciphertext too short")
	}

	// CBC mode always works in whole blocks.
	if len(ct)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	bs := block.BlockSize()
	blocks := len(ct) / bs

	for i := 0; i < blocks; i++ {
		block.Decrypt(ct[bs*i:bs*(i+1)], ct[bs*i:bs*(i+1)])
	}

	return ct, nil

}

func ECBDetectOracle(data []byte) bool {
	l := len(data)
	m := make(map[string]int)
	for i := 0; i < l/16; i++ {
		s := string(data[i*16 : (i+1)*16])
		if m[s] == 1 {
			return true
		}
		m[s] = 1
	}

	return false
}
