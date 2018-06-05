package transaction

import (
	"bytes"
	"testing"

	"github.com/mileschao/echain/common"
	"github.com/ontio/ontology-crypto/keypair"
)

func TestTxSerialize(t *testing.T) {
	var tx Transaction
	tx.Version = 0xFF
	tx.TxType = Deploy
	tx.Nonce = 0xFFFF
	tx.GasPrice = 0x01
	tx.GasLimit = 10000
	tx.Payer = common.Address{
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
		0xFF, 0xFE, 0xFD, 0xFC, 0xFB,
	}
	tx.Payload = nil
	var txAttr = &TxAttribute{
		Usage: Nonce,
		Data:  []byte{0xFF},
		Size:  uint32(1),
	}
	tx.Attributes = make([]*TxAttribute, 0)
	tx.Attributes = append(tx.Attributes, txAttr)
	var sig Sig
	_, pk, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	if err != nil {
		t.Errorf("generate public key:%s", err)
	}
	sig.PubKeys = []keypair.PublicKey{pk}
	sig.M = 8
	sig.SigData = [][]byte{[]byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB}}
	tx.Sigs = make([]*Sig, 0)
	tx.Sigs = append(tx.Sigs, &sig)

	buf := new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		t.Errorf("tx serialize: %s", err)
	}

	var tx2 Transaction
	if err := tx2.Deserialize(buf); err != nil {
		t.Errorf("tx deserialize: %s", err)
	}

	if tx2.Sigs[0].M != tx.Sigs[0].M {
		t.Errorf("tx deserialize")
	}

}
