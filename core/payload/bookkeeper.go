package payload

import (
	"encoding/binary"
	"io"

	"github.com/mileschao/echain/common/serialize"
	"github.com/ontio/ontology-crypto/keypair"
)

//BookkeeperAction bookkeeper action type
type BookkeeperAction byte

const (
	// BookkeeperActionADD action add bookkeeper
	BookkeeperActionADD BookkeeperAction = 0
	// BookkeeperActionSUB action sub bookkeeper
	BookkeeperActionSUB BookkeeperAction = 1
)

// Bookkeeper is an implementation of transaction payload for consensus bookkeeper list modification
type Bookkeeper struct {
	PubKey keypair.PublicKey
	Action BookkeeperAction
	Cert   []byte
	Issuer keypair.PublicKey
}

// Serialize implement Payload interface
func (bk *Bookkeeper) Serialize(w io.Writer) error {
	pk := keypair.SerializePublicKey(bk.PubKey)
	var bkvb = &serialize.VarBytes{
		Len:   uint64(len(pk)),
		Bytes: pk,
	}
	if err := bkvb.Serialize(w); err != nil {
		return err
	}
	binary.Write(w, binary.LittleEndian, bk.Action)
	var certvb = &serialize.VarBytes{
		Len:   uint64(len(bk.Cert)),
		Bytes: bk.Cert,
	}
	if err := certvb.Serialize(w); err != nil {
		return err
	}
	issuer := keypair.SerializePublicKey(bk.Issuer)
	var issuservb = &serialize.VarBytes{
		Len:   uint64(len(issuer)),
		Bytes: issuer,
	}
	return issuservb.Serialize(w)
}

// Deserialize deserialize Bookkeeper from io.Reader
func (bk *Bookkeeper) Deserialize(r io.Reader) error {
	var pkvb serialize.VarBytes
	if err := pkvb.Deserialize(r); err != nil {
		return err
	}
	var err error
	bk.PubKey, err = keypair.DeserializePublicKey(pkvb.Bytes)
	if err != nil {
		return err
	}
	binary.Read(r, binary.LittleEndian, &bk.Action)
	var certvb serialize.VarBytes
	if err := certvb.Deserialize(r); err != nil {
		return err
	}
	bk.Cert = certvb.Bytes
	var issuevb serialize.VarBytes
	if err := issuevb.Deserialize(r); err != nil {
		return err
	}
	bk.Issuer, err = keypair.DeserializePublicKey(issuevb.Bytes)
	return err
}
