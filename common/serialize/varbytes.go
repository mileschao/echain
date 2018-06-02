package serialize

import (
	"encoding/binary"
	"io"
)

// VarBytes variable bytes array with length
// serialized with format below:
// length + bytes[]
type VarBytes struct {
	Len   uint64
	Bytes []byte
}

// Serialize implement Serializable interface
// serialize a variable bytes array into format below:
// variable length + bytes[]
func (vb *VarBytes) Serialize(w io.Writer) error {
	var varlen = VarUint{UintType: GetUintTypeByValue(vb.Len), Value: vb.Len}
	if err := varlen.Serialize(w); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, vb.Bytes)
}

// Deserialize implement Deserialiazable interface
// deserialize a variable bytes arrary from buffer
// see Serialize above as reference
func (vb *VarBytes) Deserialize(r io.Reader) error {
	var varlen VarUint
	if err := varlen.Deserialize(r); err != nil {
		return err
	}
	vb.Len = uint64(varlen.Value)
	vb.Bytes = make([]byte, vb.Len)
	return binary.Read(r, binary.LittleEndian, vb.Bytes)
}
