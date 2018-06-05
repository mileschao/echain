package payload

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common/serialize"

	"github.com/mileschao/echain/smartcontract/types"
)

//DeployCode deploy code payload
type DeployCode struct {
	Code        types.VMCode
	NeedStorage bool
	Name        string
	Version     string
	Author      string
	Email       string
	Description string
}

// Serialize implement Payload interface
func (dc *DeployCode) Serialize(w io.Writer) error {

	if err := dc.Code.Serialize(w); err != nil {
		return err
	}
	binary.Write(w, binary.LittleEndian, dc.NeedStorage)
	var nvb = &serialize.VarBytes{
		Len:   uint64(len(dc.Name)),
		Bytes: []byte(dc.Name),
	}
	if err := nvb.Serialize(w); err != nil {
		return err
	}
	var vervb = &serialize.VarBytes{
		Len:   uint64(len(dc.Version)),
		Bytes: []byte(dc.Version),
	}
	if err := vervb.Serialize(w); err != nil {
		return err
	}
	var authorvb = &serialize.VarBytes{
		Len:   uint64(len(dc.Author)),
		Bytes: []byte(dc.Author),
	}
	if err := authorvb.Serialize(w); err != nil {
		return err
	}
	var emailvb = &serialize.VarBytes{
		Len:   uint64(len(dc.Email)),
		Bytes: []byte(dc.Email),
	}
	if err := emailvb.Serialize(w); err != nil {
		return err
	}
	var descvb = &serialize.VarBytes{
		Len:   uint64(len(dc.Description)),
		Bytes: []byte(dc.Description),
	}
	return descvb.Serialize(w)
}

//Deserialize implement Payload interface
func (dc *DeployCode) Deserialize(r io.Reader) error {
	if err := dc.Code.Deserialize(r); err != nil {
		return err
	}
	binary.Read(r, binary.LittleEndian, &dc.NeedStorage)
	var nvb serialize.VarBytes
	if err := nvb.Deserialize(r); err != nil {
		return err
	}
	dc.Name = string(nvb.Bytes)
	var vervb serialize.VarBytes
	if err := vervb.Deserialize(r); err != nil {
		return err
	}
	dc.Version = string(vervb.Bytes)
	var authorvb serialize.VarBytes
	if err := authorvb.Deserialize(r); err != nil {
		return err
	}
	dc.Author = string(authorvb.Bytes)
	var emailvb serialize.VarBytes
	if err := emailvb.Deserialize(r); err != nil {
		return err
	}
	dc.Email = string(emailvb.Bytes)
	var descvb serialize.VarBytes
	if err := descvb.Deserialize(r); err != nil {
		return err
	}
	dc.Description = string(descvb.Bytes)
	return nil
}

//Bytes get byte array
func (dc *DeployCode) Bytes() []byte {
	b := new(bytes.Buffer)
	dc.Serialize(b)
	return b.Bytes()
}
