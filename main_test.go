package main

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"
)

func TestHexToBase64(t *testing.T) {
	hex := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	base64 := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	res, err := HexToBase64(hex)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res != base64 {
		t.Errorf("Base64 result does not match expected: \n%s vs \n%s", res, base64)
	}
}

func TestXOR(t *testing.T) {
	hex1 := "1c0111001f010100061a024b53535009181c"
	hex2 := "686974207468652062756c6c277320657965"
	hex_expected := "746865206b696420646f6e277420706c6179"
	b1, err := hex.DecodeString(hex1)
	if err != nil {
		t.Errorf(err.Error())
	}
	b2, err := hex.DecodeString(hex2)
	if err != nil {
		t.Errorf(err.Error())
	}

	res, err := XOR(b1, b2)
	if err != nil {
		t.Errorf(err.Error())
	}
	hex_res := hex.EncodeToString(res)
	if hex_res != hex_expected {
		t.Errorf("Did not get expected result from XOR: \n%s vs. \n%s", hex_res, hex_expected)
	}
}

func TestSingleCharXOR(t *testing.T) {
	ct := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	expected := "7a56565250575e19747a1e4a195550525c19581949564c575d19565f195b585a5657"
	b_ct, err := hex.DecodeString(ct)
	if err != nil {
		t.Errorf(err.Error())
	}

	res := RepeatingXOR(b_ct, []byte("a"))
	hex_res := hex.EncodeToString(res)
	if hex_res != expected {
		t.Errorf("Did not get expected result from SingleCharXOR with 'a': \n%s vs. \n%s", hex_res, expected)
	}
}

func TestSingleCharOracle(t *testing.T) {
	ct := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	b_ct, _ := hex.DecodeString(ct)
	res := SingleCharOracle(b_ct)
	if string(res) != "X" {
		t.Errorf("Did not get expected result from SingleCharOracle: \n%s vs. \n%s", res, "X")
	}
	pt := RepeatingXOR(b_ct, []byte(res))
	expected := "Cooking MC's like a pound of bacon"
	if string(pt) != expected {
		t.Errorf("Did not get expected result from SingleCharOracle: \n%s vs. \n%s", pt, expected)
	}
}

func TestRepeatingKeyOracle(t *testing.T) {
	file, err := os.Open("challenge6.txt") // For read access.
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestHamming(t *testing.T) {
	buf1 := []byte("this is a test")
	buf2 := []byte("wokka wokka!!!")
	dis := Hamming(buf1, buf2)
	if dis != 37 {
		t.Errorf("Expected hamming distance is 37. Got %d", dis)
	}
}
