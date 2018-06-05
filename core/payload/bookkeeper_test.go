package payload

import (
	"bytes"
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
)

func TestBookkeeperSerialize(t *testing.T) {
	_, bkpk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	if err != nil {
		t.Errorf("generate key: %s", err)
	}
	_, ispk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	if err != nil {
		t.Errorf("generate key: %s", err)
	}
	var bk = &Bookkeeper{
		PubKey: bkpk,
		Action: BookkeeperActionADD,
		Cert:   []byte{0xFF},
		Issuer: ispk,
	}
	buf := new(bytes.Buffer)
	if err := bk.Serialize(buf); err != nil {
		t.Errorf("book keeper serialize: %s", err)
	}

	var bk2 Bookkeeper
	if err := bk2.Deserialize(buf); err != nil {
		t.Errorf("book keeper deserialize: %s", err)
	}
	if !bytes.Equal(bk.Cert, bk2.Cert) {
		t.Errorf("book keeper deserialize")
	}
}
