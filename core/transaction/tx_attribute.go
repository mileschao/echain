package transaction

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"

	"github.com/mileschao/echain/common/serialize"
)

var (
	//ErrUnSupportUsageType unsupport usage type
	ErrUnSupportUsageType = errors.New("unsupport usage type")
)

// TxAttrUsage transaction attribute usage type
// Nonce,Script,DesciptionUrl, Description
type TxAttrUsage byte

const (
	// Nonce TxAttrUsage
	Nonce TxAttrUsage = 0x00
	// Script TxAttrUsage
	Script TxAttrUsage = 0x20
	// DescriptionURL TxAttrUsage
	DescriptionURL TxAttrUsage = 0x81
	// Description TxAttrUsage
	Description TxAttrUsage = 0x90
)

func (txattr TxAttrUsage) isValid() bool {
	if txattr != Nonce &&
		txattr != Script &&
		txattr != DescriptionURL &&
		txattr != Description {
		return false
	}
	return true
}

// TxAttribute transaction attribute
type TxAttribute struct {
	Usage TxAttrUsage // Nonce, Script, DescriptionURL, Description
	Data  []byte
	Size  uint32
}

// NewTxAttribute create an new transaction attribute
func NewTxAttribute(usage TxAttrUsage, data []byte) TxAttribute {
	tx := TxAttribute{
		Usage: usage,
		Data:  data,
	}
	if usage == DescriptionURL {
		us := reflect.TypeOf(tx.Usage).Size()
		ss := reflect.TypeOf(tx.Usage).Size()
		tx.Size = uint32(us) + uint32(ss) + uint32(len(tx.Data))
	} else {
		tx.Size = 0
	}
	return tx
}

//Serialize implement Serializable interface
func (tx *TxAttribute) Serialize(w io.Writer) error {

	if err := binary.Write(w, binary.LittleEndian, byte(tx.Usage)); err != nil {
		return err
	}
	if !tx.Usage.isValid() {
		return ErrUnSupportUsageType
	}
	var vbd = &serialize.VarBytes{
		Len:   uint64(len(tx.Data)),
		Bytes: tx.Data,
	}

	return vbd.Serialize(w)
}

//Deserialize implement Serializable interface
func (tx *TxAttribute) Deserialize(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &tx.Usage)
	if !tx.Usage.isValid() {
		return ErrUnSupportUsageType
	}
	var vbd serialize.VarBytes
	return vbd.Deserialize(r)
}

//Bytes get the byte array of TxAttribute
func (tx *TxAttribute) Bytes() []byte {
	bf := new(bytes.Buffer)
	tx.Serialize(bf)
	return bf.Bytes()
}
