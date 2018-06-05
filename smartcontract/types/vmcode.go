package types

import (
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common"
	"github.com/mileschao/echain/common/serialize"
	"golang.org/x/crypto/ripemd160"
)

//VMType virtual machine type
type VMType byte

const (
	//Native native virtual machine
	Native = VMType(0xFF)
	//NEOVM NEO virtual machine
	NEOVM = VMType(0x80)
	//WASMVM WASM virtual machine
	WASMVM = VMType(0x90)
	// EVM = VmType(0x90)
)

// VMCode describe smart contract code and vm type
type VMCode struct {
	VMType VMType
	Code   []byte
}

// Serialize implement serilaizable interface
func (vc *VMCode) Serialize(w io.Writer) error {
	binary.Write(w, binary.LittleEndian, vc.VMType)
	var vcvb = &serialize.VarBytes{
		Len:   uint64(len(vc.Code)),
		Bytes: vc.Code,
	}
	return vcvb.Serialize(w)
}

// Deserialize implement serializable interface
func (vc *VMCode) Deserialize(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &vc.VMType)
	var vcvb serialize.VarBytes
	if err := vcvb.Deserialize(r); err != nil {
		return err
	}
	vc.Code = vcvb.Bytes
	return nil
}

// Address return address of contract
func (vc *VMCode) Address() common.Address {
	var addr common.Address
	temp := sha256.Sum256(vc.Code)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])

	addr[0] = byte(vc.VMType)
	return addr
}

// IsVMCodeAddress check whether address is smart contract address
func IsVMCodeAddress(addr common.Address) bool {
	vmType := addr[0]
	if vmType == byte(Native) || vmType == byte(NEOVM) || vmType == byte(WASMVM) {
		return true
	}

	return false
}
