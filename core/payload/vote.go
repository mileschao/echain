package payload

import (
	"io"

	"github.com/mileschao/echain/common/serialize"

	"github.com/mileschao/echain/common"
	"github.com/ontio/ontology-crypto/keypair"
)

const (
	// MaxVoteKeys max votation public keys
	MaxVoteKeys = 1024
)

//Vote votation payload
type Vote struct {
	// vote public keys list
	PubKeys []keypair.PublicKey
	Account common.Address
}

//isValid
func (v *Vote) isValid() bool {
	if len(v.PubKeys) > MaxVoteKeys {
		return false
	}
	return true
}

// Serialize implement Payload interface
func (v *Vote) Serialize(w io.Writer) error {
	var pkvu = &serialize.VarUint{
		UintType: serialize.GetUintTypeByValue(uint64(len(v.PubKeys))),
		Value:    uint64(len(v.PubKeys)),
	}
	if err := pkvu.Serialize(w); err != nil {
		return err
	}
	for _, pk := range v.PubKeys {
		buf := keypair.SerializePublicKey(pk)
		var pkvb = &serialize.VarBytes{
			Len:   uint64(len(buf)),
			Bytes: buf,
		}
		if err := pkvb.Serialize(w); err != nil {
			return err
		}
	}
	return v.Account.Serialize(w)
}

//Deserialize implement Payload interface
func (v *Vote) Deserialize(r io.Reader) error {
	var pkvu serialize.VarUint
	if err := pkvu.Deserialize(r); err != nil {
		return err
	}
	v.PubKeys = make([]keypair.PublicKey, 0)
	for i := uint64(0); i < pkvu.Value; i++ {
		var pkvb serialize.VarBytes
		if err := pkvb.Deserialize(r); err != nil {
			return err
		}
		pk, err := keypair.DeserializePublicKey(pkvb.Bytes)
		if err != nil {
			return err
		}
		v.PubKeys = append(v.PubKeys, pk)
	}
	return v.Account.Deserialize(r)
}
