package transaction

import (
	"bytes"
	"testing"

	"github.com/ontio/ontology-crypto/keypair"
)

func TestSignatureSerialize(t *testing.T) {
	var sig Sig
	_, pk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	if err != nil {
		t.Errorf("generate public key:%s", err)
	}
	sig.PubKeys = []keypair.PublicKey{pk}
	sig.M = 8
	sig.SigData = [][]byte{[]byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB}}

	buf := new(bytes.Buffer)
	if err := sig.Serialize(buf); err != nil {
		t.Errorf("serialize: %s", err)
	}
	var sig2 Sig
	if err := sig2.Deserialize(buf); err != nil {
		t.Errorf("deserialize: %s", err)
	}
	if sig.M != sig2.M {
		t.Errorf("deserialize")
	}
}
