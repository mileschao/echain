package payload

import (
	"bytes"
	"testing"

	"github.com/mileschao/echain/common"
	"github.com/ontio/ontology-crypto/keypair"
)

func TestVoteSerialize(t *testing.T) {
	_, pk, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	var vote = &Vote{
		PubKeys: []keypair.PublicKey{pk},
		Account: common.ADDRESS_EMPTY,
	}
	buf := new(bytes.Buffer)
	if err := vote.Serialize(buf); err != nil {
		t.Errorf("vote serialize: %s", err)
	}
	var vote2 Vote
	if err := vote2.Deserialize(buf); err != nil {
		t.Errorf("vote deserialize: %s", err)
	}
	if !keypair.ComparePublicKey(vote.PubKeys[0], vote2.PubKeys[0]) {
		t.Errorf("vote deserialize")
	}

}
