package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

// helper function t convert given value to a hash one
func hashFunc(val, alg string) string {
	var h hash.Hash
	if alg == "sha256" {
		h = sha256.New()
	} else if alg == "sha512" {
		h = sha512.New()
	} else {
		h = sha1.New()
	}
	h.Write([]byte(val))
	return hex.EncodeToString(h.Sum(nil))
}
