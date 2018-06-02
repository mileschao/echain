package common

import (
	"encoding/binary"
	"errors"
	"io"
)

/*
 * uint256: unsigned integer of 256 bit size
 * 32 * 8 = 256
 * 1 byte = 8 bit
 * actually this is base256
 */

// UINT256_SIZE 32 bytes
const UINT256_SIZE = 32

// Uint256 base256 with fix size byte array format
type Uint256 [UINT256_SIZE]byte

var (
	// UINT256_EMPTY empty uint256
	UINT256_EMPTY = Uint256{}

	// ErrBytesSize the byte array must be 32
	ErrBytesSize = errors.New("wrong bytes array size")
)

// Serialize implement Serializable interface
// serialize fix size byte array
func (u *Uint256) Serialize(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, u)
}

// Deserialize implement Serializable interface
// deserialize buffer to fix size byte array
func (u *Uint256) Deserialize(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, u)
}

// Bytes return bytes with copied content
func (u *Uint256) Bytes() []byte {
	b := make([]byte, UINT256_SIZE)
	copy(b, u[:])
	return b
}

// FromBytes copy content from bytes array to uint256
func (u *Uint256) FromBytes(b []byte) error {
	if len(b) != UINT256_SIZE {
		return ErrBytesSize
	}
	copy(u[:], b[:])
	return nil
}
