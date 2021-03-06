package common

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUint256Serialize(t *testing.T) {
	b := new(bytes.Buffer)
	var u256 = Uint256{0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD}
	if err := u256.Serialize(b); err != nil {
		t.Errorf("uint256 seriliaze: %s", err)
	}
	bs := b.Bytes()
	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 serialize: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256Deserialize(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD})
	var u256 = UINT256_EMPTY
	bs := b.Bytes()
	if err := u256.Deserialize(b); err != nil {
		t.Errorf("uint256 deserialize: %s", err)
	}

	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 serialize: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256Bytes(t *testing.T) {
	var u256 = Uint256{0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD}
	bs := u256.Bytes()
	if !reflect.DeepEqual(bs[:], u256[:]) {
		t.Errorf("uint256 get bytes: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}

func TestUint256FromBytes(t *testing.T) {
	b := bytes.NewBuffer([]byte{
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF,
		0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE,
		0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD,
		0xFE, 0xFF, 0xFD, 0xFE, 0xFF, 0xFD, 0xFD, 0xFD})
	var u256 = UINT256_EMPTY
	bs := b.Bytes()
	if err := u256.FromBytes(bs); err != nil {
		t.Errorf("uint256 frombytes: %s", err)
	}

	if !reflect.DeepEqual(u256[:], bs[:]) {
		t.Errorf("uint256 frombytes: \nuint256-%X\nbuffer -%X", u256, bs)
	}
}
