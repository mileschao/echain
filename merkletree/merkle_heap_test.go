package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/mileschao/echain/common"
)

func TestProof(t *testing.T) {

	var u1 common.Uint256
	var u2 common.Uint256
	var u3 common.Uint256
	var u4 common.Uint256
	var u5 common.Uint256
	var u6 common.Uint256
	var u7 common.Uint256
	var u8 common.Uint256
	var u9 common.Uint256
	var u10 common.Uint256
	var u11 common.Uint256

	b1, _ := hex.DecodeString("f3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l1 = sha256.Sum256(b1)
	b2, _ := hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l2 = sha256.Sum256(b2)
	b3, _ := hex.DecodeString("d3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l3 = sha256.Sum256(b3)
	b4, _ := hex.DecodeString("c3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l4 = sha256.Sum256(b4)
	b5, _ := hex.DecodeString("b3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l5 = sha256.Sum256(b5)
	b6, _ := hex.DecodeString("a3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l6 = sha256.Sum256(b6)
	b7, _ := hex.DecodeString("93b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l7 = sha256.Sum256(b7)
	b8, _ := hex.DecodeString("83b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l8 = sha256.Sum256(b8)
	b9, _ := hex.DecodeString("73b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l9 = sha256.Sum256(b9)
	b10, _ := hex.DecodeString("63b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l10 = sha256.Sum256(b10)
	b11, _ := hex.DecodeString("53b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	var l11 = sha256.Sum256(b11)

	u1.FromBytes(l1[:])
	u2.FromBytes(l2[:])
	u3.FromBytes(l3[:])
	u4.FromBytes(l4[:])
	u5.FromBytes(l5[:])
	u6.FromBytes(l6[:])
	u7.FromBytes(l7[:])
	u8.FromBytes(l8[:])
	u9.FromBytes(l9[:])
	u10.FromBytes(l10[:])
	u11.FromBytes(l11[:])

	var hu0 = []common.Uint256{}
	var hm = &memoryHashStorage{
		hashes: make([]common.Uint256, 0, 19),
	}
	mh := NewMerkleStorage(0, hu0, hm)
	mh.AddLeaf(u1)
	mh.AddLeaf(u2)
	mh.AddLeaf(u3)
	mh.AddLeaf(u4)
	mh.AddLeaf(u5)
	mh.AddLeaf(u6)
	mh.AddLeaf(u7)
	mh.AddLeaf(u8)
	mh.AddLeaf(u9)
	mh.AddLeaf(u10)
	mh.AddLeaf(u11)

	if mh.LeafSize() != 11 {
		t.Errorf("mh add leaf error: %d", mh.LeafSize())
	}
	u1s, _ := mh.hashStorage.GetHash(0)
	if !reflect.DeepEqual(u1s, u1) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 1, u1s.Bytes(), u1.Bytes())
	}

	u2s, _ := mh.hashStorage.GetHash(1)
	if !reflect.DeepEqual(u2s, u2) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 2, u2s.Bytes(), u2.Bytes())
	}

	h12 := nodeHash(u1, u2)
	u12s, _ := mh.hashStorage.GetHash(2)
	if !reflect.DeepEqual(u12s, h12) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 3, u12s.Bytes(), h12.Bytes())
	}

	u3s, _ := mh.hashStorage.GetHash(3)
	if !reflect.DeepEqual(u3s, u3) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 4, u3s.Bytes(), u3.Bytes())
	}

	u4s, _ := mh.hashStorage.GetHash(4)
	if !reflect.DeepEqual(u4s, u4) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 5, u4s.Bytes(), u4.Bytes())
	}

	h34 := nodeHash(u3, u4)
	u34s, _ := mh.hashStorage.GetHash(5)
	if !reflect.DeepEqual(u34s, h34) {
		t.Errorf("merkle heap add leaf: %d\n%X\n%X", 6, u34s.Bytes(), h34.Bytes())
	}

	proofs, err := mh.InclusionProof(4, 11)
	if err != nil {
		t.Errorf("error proof: %s", err)
	}
	for i, p := range proofs {
		fmt.Printf("proof: i: %d, %X\n", i, p)
	}
}
