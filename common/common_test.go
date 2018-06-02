package common

import (
	"reflect"
	"strings"
	"testing"
)

func TestNonce(t *testing.T) {
	r1 := Nonce()
	r2 := Nonce()
	if r1 == r2 {
		t.Errorf("nonce: %d, %d", r1, r2)
	}
}

func TestHex(t *testing.T) {
	b := []byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA}
	h := Hex(b)
	if strings.ToLower("FFFEFDFCFBFA") != h {
		t.Errorf("hex: %s", h)
	}
}

func TestHexToBytes(t *testing.T) {
	b := []byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA}
	h := Hex(b)
	bs, err := HexToBytes(h)
	if err != nil {
		t.Errorf("hex to bytes: %s", err)
	}
	if !reflect.DeepEqual(bs[:], b[:]) {
		t.Errorf("hex to bytes:\n0-%X,\n1-%X", b[:], bs[:])
	}
}

func TestFileExist(t *testing.T) {
	if FileExisted("dumy/foo.txt") {
		t.Errorf("common fileexits")
	}
}
