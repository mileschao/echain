package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common"
	"github.com/mileschao/echain/common/serialize"
	"github.com/mileschao/echain/core/payload"
	stypes "github.com/mileschao/echain/smartcontract/types"
)

//TxType transaction type
type TxType byte

const (
	//Bookkeeper transaction type
	Bookkeeper TxType = 0x02
	//Claim transaction type
	Claim TxType = 0x03
	//Deploy transaction type
	Deploy TxType = 0xd0
	//Invoke transaction type
	Invoke TxType = 0xd1
	//Enrollment transaction type
	Enrollment TxType = 0x04
	//Vote transaction type
	Vote TxType = 0x05
)

// Transaction transaction
type Transaction struct {
	Version    byte
	TxType     TxType
	Nonce      uint32
	GasPrice   uint64
	GasLimit   uint64
	Payer      common.Address
	Payload    payload.Payload
	Attributes []*TxAttribute
	Sigs       []*Sig
	hash       *common.Uint256
}

// NewDeployTx returns a deploy Transaction
func NewDeployTx(code stypes.VMCode, name, version, author, email, desp string, needStorage bool) *Transaction {
	//TODO: check arguments
	DeployCodePayload := &payload.DeployCode{
		Code:        code,
		NeedStorage: needStorage,
		Name:        name,
		Version:     version,
		Author:      author,
		Email:       email,
		Description: desp,
	}

	return &Transaction{
		TxType:     Deploy,
		Payload:    DeployCodePayload,
		Attributes: nil,
	}
}

// NewInvokeTx returns an invoke Transaction
func NewInvokeTx(vmcode stypes.VMCode) *Transaction {
	//TODO: check arguments
	invokeCodePayload := &payload.InvokeCode{
		Code: vmcode,
	}

	return &Transaction{
		TxType:     Invoke,
		Payload:    invokeCodePayload,
		Attributes: nil,
	}
}

//Serialize implement Payload interface
func (tx *Transaction) Serialize(w io.Writer) error {
	binary.Write(w, binary.LittleEndian, tx.Version)
	binary.Write(w, binary.LittleEndian, byte(tx.TxType))
	binary.Write(w, binary.LittleEndian, tx.Nonce)
	binary.Write(w, binary.LittleEndian, tx.GasPrice)
	binary.Write(w, binary.LittleEndian, tx.GasLimit)
	if err := tx.Payer.Serialize(w); err != nil {
		return err
	}
	if err := tx.Payload.Serialize(w); err != nil {
		return err
	}

	var attrvu = &serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(tx.Attributes))),
		Value:    uint64(len(tx.Attributes)),
	}
	if err := attrvu.Serialize(w); err != nil {
		return err
	}
	for _, attr := range tx.Attributes {
		if err := attr.Serialize(w); err != nil {
			return err
		}
	}
	var svu = &serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(tx.Sigs))),
		Value:    uint64(len(tx.Sigs)),
	}
	if err := svu.Serialize(w); err != nil {
		return err
	}
	for _, sig := range tx.Sigs {
		if err := sig.Serialize(w); err != nil {
			return err
		}
	}
	return nil
}

// Deserialize implement the Payload interface
func (tx *Transaction) Deserialize(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &tx.Version)
	binary.Read(r, binary.LittleEndian, &tx.TxType)
	binary.Read(r, binary.LittleEndian, &tx.Nonce)
	binary.Read(r, binary.LittleEndian, &tx.GasPrice)
	binary.Read(r, binary.LittleEndian, &tx.GasLimit)
	if err := tx.Payer.Deserialize(r); err != nil {
		return err
	}
	if err := tx.Payload.Deserialize(r); err != nil {
		return err
	}
	var attrvu serialize.VarUint
	if err := attrvu.Deserialize(r); err != nil {
		return err
	}
	tx.Attributes = make([]*TxAttribute, 0, attrvu.Value)
	for i := uint64(0); i < attrvu.Value; i++ {
		var attr TxAttribute
		if err := attr.Deserialize(r); err != nil {
			return err
		}
		tx.Attributes = append(tx.Attributes, &attr)
	}

	var sigvu serialize.VarUint
	if err := attrvu.Deserialize(r); err != nil {
		return err
	}
	tx.Sigs = make([]*Sig, 0, sigvu.Value)
	for i := uint64(0); i < sigvu.Value; i++ {
		var sig Sig
		if err := sig.Deserialize(r); err != nil {
			return err
		}
		tx.Sigs = append(tx.Sigs, &sig)
	}
	return nil
}

// Hash get the transaction's hash value
func (tx *Transaction) Hash() common.Uint256 {
	if tx.hash == nil {
		w := new(bytes.Buffer)
		binary.Write(w, binary.LittleEndian, tx.Version)
		binary.Write(w, binary.LittleEndian, byte(tx.TxType))
		binary.Write(w, binary.LittleEndian, tx.Nonce)
		binary.Write(w, binary.LittleEndian, tx.GasPrice)
		binary.Write(w, binary.LittleEndian, tx.GasLimit)
		if err := tx.Payer.Serialize(w); err != nil {
			return common.UINT256_EMPTY
		}
		if err := tx.Payload.Serialize(w); err != nil {
			return common.UINT256_EMPTY
		}
		var attrvu = &serialize.VarUint{
			UintType: serialize.GetUintTypeByValue(uint64(len(tx.Attributes))),
			Value:    uint64(len(tx.Attributes)),
		}
		if err := attrvu.Serialize(w); err != nil {
			return common.UINT256_EMPTY
		}
		for _, attr := range tx.Attributes {
			if err := attr.Serialize(w); err != nil {
				return common.UINT256_EMPTY
			}
		}

		temp := sha256.Sum256(w.Bytes())
		f := common.Uint256(sha256.Sum256(temp[:]))
		tx.hash = &f
	}
	return *tx.hash
}

// SetHash set hash value of the transaction
func (tx *Transaction) SetHash(hash common.Uint256) {
	tx.hash = &hash
}

// Bytes get byte array of transaction
func (tx *Transaction) Bytes() []byte {
	b := new(bytes.Buffer)
	tx.Serialize(b)
	return b.Bytes()
}
