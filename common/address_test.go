package common

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestAddressSerialize(t *testing.T) {
	b := new(bytes.Buffer)
	var addr = Address{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB}
	if err := addr.Serialize(b); err != nil {
		t.Errorf("address serialize: %s", err)
	}
	if !reflect.DeepEqual(addr[:], b.Bytes()[:]) {
		t.Errorf("address serialize:\n 0-%X\n1-%X", addr[:], b.Bytes()[:])
	}
}

func TestAddressDeserialize(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB})
	var addr = ADDRESS_EMPTY
	bs := b.Bytes()
	if err := addr.Deserialize(b); err != nil {
		t.Errorf("address deserialize: %s", err)
	}

	if !reflect.DeepEqual(addr[:], bs[:]) {
		t.Errorf("address deserialize:\n 0-%X\n1-%X", addr[:], bs[:])
	}
}

func TestAddressFromBytes(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB})
	var addr = ADDRESS_EMPTY
	bs := b.Bytes()
	if err := addr.FromBytes(bs); err != nil {
		t.Errorf("address from bytes: %s", err)
	}
	if !reflect.DeepEqual(addr[:], bs[:]) {
		t.Errorf("address from bytes:\n0-%X\n1-%X", addr, bs)
	}
}

func TestAddressBase58(t *testing.T) {
	var addr = Address{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB}
	base58addr := addr.Base58()
	fmt.Println(base58addr)
	var addr58 = ADDRESS_EMPTY
	if err := addr58.FromBase58(base58addr); err != nil {
		t.Errorf("address base58: %s", err)
	}
}
