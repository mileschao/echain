package common

import (
	"encoding/hex"
	"math/rand"
	"os"
)

// Nonce returns random nonce
func Nonce() uint64 {
	// TODO: replace with the real random number generator
	nonce := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return nonce
}

// Hex get hex string corresponding to byte array
func Hex(data []byte) string {
	return hex.EncodeToString(data)
}

// HexToBytes convert hex string to byte array
func HexToBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

// FileExisted checks whether filename exists in filesystem
func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
