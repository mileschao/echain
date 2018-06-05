package transaction

import (
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common/serialize"
	"github.com/ontio/ontology-crypto/keypair"
)

// Sig signature
type Sig struct {
	PubKeys []keypair.PublicKey
	M       uint8
	SigData [][]byte
}

//Serialize implement Payload interface
func (s *Sig) Serialize(w io.Writer) error {
	var pkvu = &serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(s.PubKeys))),
		Value:    uint64(len(s.PubKeys)),
	}
	if err := pkvu.Serialize(w); err != nil {
		return err
	}

	for _, pk := range s.PubKeys {
		pkb := keypair.SerializePublicKey(pk)
		var pkvb = serialize.VarBytes{
			Len:   uint64(len(pkb)),
			Bytes: pkb,
		}
		pkvb.Serialize(w)
	}

	binary.Write(w, binary.LittleEndian, s.M)

	var sdvu = &serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(s.SigData))),
		Value:    uint64(len(s.SigData)),
	}

	if err := sdvu.Serialize(w); err != nil {
		return err
	}

	for _, sig := range s.SigData {
		var sgvb = &serialize.VarBytes{
			Len:   uint64(len(sig)),
			Bytes: sig,
		}
		if err := sgvb.Serialize(w); err != nil {
			return err
		}
	}

	return nil
}

//Deserialize implement Payload interface
func (s *Sig) Deserialize(r io.Reader) error {
	var pkvu serialize.VarUint
	if err := pkvu.Deserialize(r); err != nil {
		return err
	}
	s.PubKeys = make([]keypair.PublicKey, pkvu.Value)
	for i := uint64(0); i < pkvu.Value; i++ {
		var pkvb serialize.VarBytes
		if err := pkvb.Deserialize(r); err != nil {
			return err
		}
		var err error
		s.PubKeys[i], err = keypair.DeserializePublicKey(pkvb.Bytes)
		if err != nil {
			return err
		}

	}

	binary.Read(r, binary.LittleEndian, &s.M)

	var sgvu serialize.VarUint
	if err := sgvu.Deserialize(r); err != nil {
		return err
	}
	s.SigData = make([][]byte, sgvu.Value)
	for i := uint64(0); i < sgvu.Value; i++ {
		var sgvb serialize.VarBytes
		if err := sgvb.Deserialize(r); err != nil {
			return err
		}
		s.SigData[i] = sgvb.Bytes
	}

	return nil
}
