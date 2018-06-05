package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common/serialize"

	"github.com/mileschao/echain/common"
	"github.com/ontio/ontology-crypto/keypair"
)

// Header block header
type Header struct {
	Version          uint32
	PrevBlockHash    common.Uint256
	TransactionsRoot common.Uint256
	BlockRoot        common.Uint256
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	ConsensusPayload []byte
	NextBookkeeper   common.Address
	Bookkeepers      []keypair.PublicKey
	SigData          [][]byte
	hash             *common.Uint256
}

//Serialize implement the Serializable interface
func (bh *Header) Serialize(w io.Writer) error {
	binary.Write(w, binary.LittleEndian, bh.Version)
	bh.PrevBlockHash.Serialize(w)
	bh.TransactionsRoot.Serialize(w)
	bh.BlockRoot.Serialize(w)
	binary.Write(w, binary.LittleEndian, bh.Timestamp)
	binary.Write(w, binary.LittleEndian, bh.Height)
	binary.Write(w, binary.LittleEndian, bh.ConsensusData)
	var cpvb = serialize.VarBytes{Bytes: bh.ConsensusPayload, Len: uint64(len(bh.ConsensusPayload))}
	cpvb.Serialize(w)
	bh.NextBookkeeper.Serialize(w)

	var bkvu = serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(bh.Bookkeepers))),
		Value:    uint64(len(bh.Bookkeepers)),
	}
	bkvu.Serialize(w)

	for _, bk := range bh.Bookkeepers {
		bkb := keypair.SerializePublicKey(bk)
		var bkvb = serialize.VarBytes{
			Len:   uint64(len(bkb)),
			Bytes: bkb,
		}
		bkvb.Serialize(w)
	}

	var sigvu = serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(bh.SigData))),
		Value:    uint64(len(bh.SigData)),
	}
	sigvu.Serialize(w)

	for _, sg := range bh.SigData {
		var sgb = serialize.VarBytes{
			Len:   uint64(len(sg)),
			Bytes: sg,
		}
		sgb.Serialize(w)
	}

	return nil
}

//Deserialize implement Serializable interface
func (bh *Header) Deserialize(r io.Reader) error {
	binary.Read(r, binary.LittleEndian, &bh.Version)
	bh.PrevBlockHash.Deserialize(r)
	bh.TransactionsRoot.Deserialize(r)
	bh.BlockRoot.Deserialize(r)
	binary.Read(r, binary.LittleEndian, &bh.Timestamp)
	binary.Read(r, binary.LittleEndian, &bh.Height)
	binary.Read(r, binary.LittleEndian, &bh.ConsensusData)
	var cpvb serialize.VarBytes
	cpvb.Deserialize(r)
	bh.ConsensusPayload = cpvb.Bytes
	bh.NextBookkeeper.Deserialize(r)

	var bkvu serialize.VarUint
	bkvu.Deserialize(r)

	for i := uint64(0); i < bkvu.Value; i++ {

		var bkvb serialize.VarBytes
		bkvb.Deserialize(r)
		kp, err := keypair.DeserializePublicKey(bkvb.Bytes[:])
		if err != nil {
			// TODO: error handle
			continue
		}
		bh.Bookkeepers = append(bh.Bookkeepers, kp)
	}

	var sigvu serialize.VarUint
	sigvu.Deserialize(r)

	for i := uint64(0); i < sigvu.Value; i++ {
		var sgb serialize.VarBytes
		sgb.Deserialize(r)
		bh.SigData = append(bh.SigData, sgb.Bytes)
	}
	return nil
}

// Hash get the hash value of header
func (bh *Header) Hash() common.Uint256 {
	if bh.hash != nil {
		return *bh.hash
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, bh.Version)
	binary.Write(buf, binary.LittleEndian, bh.PrevBlockHash)
	binary.Write(buf, binary.LittleEndian, bh.TransactionsRoot)
	bh.BlockRoot.Serialize(buf)
	binary.Write(buf, binary.LittleEndian, bh.Timestamp)
	binary.Write(buf, binary.LittleEndian, bh.Height)
	binary.Write(buf, binary.LittleEndian, bh.ConsensusData)
	var cpvb = &serialize.VarBytes{
		Len:   uint64(len(bh.ConsensusPayload)),
		Bytes: bh.ConsensusPayload,
	}
	cpvb.Serialize(buf)
	bh.NextBookkeeper.Serialize(buf)

	tmp := sha256.Sum256(buf.Bytes())
	hash := common.Uint256(sha256.Sum256(tmp[:]))
	bh.hash = &hash

	return hash
}

// Bytes get header serialze byte array
func (bh *Header) Bytes() []byte {
	bf := new(bytes.Buffer)
	bh.Serialize(bf)
	return bf.Bytes()
}
