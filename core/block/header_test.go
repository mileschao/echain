package block

import (
	"bytes"
	"testing"

	"github.com/ontio/ontology-crypto/keypair"

	"github.com/mileschao/echain/common"
)

func TestHeaderSerialze(t *testing.T) {
	/*//  9393ee15ce6612484ab5be3bbc78c82af8dc0e07
	//179393EE15CE6612484AB5BE3BBC78C82AF8DC0E0778A49160
	v := "00c2eb0b"
	vb, _ := hex.DecodeString(v)
	a := neotypes.ConvertBytesToBigInteger(vb)
	fmt.Println(a.String())

	bf, _ := hex.DecodeString("179393EE15CE6612484AB5BE3BBC78C82AF8DC0E0778A49160")
	bi := new(big.Int).SetBytes(bf).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	fmt.Printf("%s\n", string(encoded))

	bf2, _ := hex.DecodeString("9393ee15ce6612484ab5be3bbc78c82af8dc0e07")
	bi2 := neotypes.ConvertBytesToBigInteger(bf2)
	bs := bi2.String()
	encoded2, _ := base58.BitcoinEncoding.Encode([]byte(bs))
	fmt.Printf("%s\n", string(encoded2))
	*/

	var head = &Header{
		Version:          0xFE,
		PrevBlockHash:    common.UINT256_EMPTY,
		TransactionsRoot: common.UINT256_EMPTY,
		BlockRoot:        common.UINT256_EMPTY,
		Timestamp:        0xFD,
		Height:           1024,
		ConsensusData:    600837,
		ConsensusPayload: nil,
		NextBookkeeper:   common.ADDRESS_EMPTY,
		Bookkeepers:      nil,
		SigData:          nil,
		hash:             &common.UINT256_EMPTY,
	}
	head.ConsensusPayload = []byte{0xFE, 0xFF, 0xFD}
	var ab = []byte{
		0xF1, 0xF2, 0xF3, 0xF1, 0xF2,
		0xF3, 0xF1, 0xF2, 0xF3, 0xF1,
		0xF2, 0xF3, 0xF1, 0xF2, 0xF3,
		0xF1, 0xF2, 0xF3, 0xF1, 0xF2,
	}
	var addr common.Address
	if err := addr.FromBytes(ab); err != nil {
		t.Errorf("address: %s", err)
	}
	head.NextBookkeeper = addr
	_, k, err := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P224)
	if err != nil {
		t.Errorf("generate key pair: %s", err)
	}
	head.Bookkeepers = []keypair.PublicKey{k}
	head.SigData = [][]byte{ab}
	bf := new(bytes.Buffer)
	if err := head.Serialize(bf); err != nil {
		t.Errorf("header serialize: %s", err)
	}

	var head2 Header
	if err := head2.Deserialize(bf); err != nil {
		t.Errorf("header deserialze: %s", err)
	}
}
