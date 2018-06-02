package merkletree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"

	"github.com/mileschao/echain/common"
)

func TestCountBit(t *testing.T) {
	var n uint64 = 0xFE
	if countBit(n) != 7 {
		t.Errorf("count bit: %b", n)
	}
}

func TestIsPower2(t *testing.T) {
	var n1 uint64 = 4
	var n2 uint64 = 5
	if !isPower2(n1) {
		t.Errorf("is power of 2: %d", n1)
	}
	if isPower2(n2) {
		t.Errorf("is power of 2: %d", n2)
	}
}

func TestHightBit(t *testing.T) {
	var h uint64 = 0x00000010
	if highBit(h) != 5 {
		t.Errorf("highest bit: %b", h)
	}
}

func TestLowBit(t *testing.T) {
	var l uint64 = 0xF0000010
	if lowBit(l) != 5 {
		t.Errorf("lowest bit: %b", l)
	}
}

func TestHighestTreeIndex(t *testing.T) {
	if 7 != highestTreeIndex(4) {
		t.Errorf("highest tree index: %d", highestTreeIndex(5))
	}
	if 7 != highestTreeIndex(5) {
		t.Errorf("highest tree index: %d", highestTreeIndex(5))
	}
	if 7 != highestTreeIndex(6) {
		t.Errorf("highest tree index: %d", highestTreeIndex(5))
	}
	if 7 != highestTreeIndex(7) {
		t.Errorf("highest tree index: %d", highestTreeIndex(5))
	}
	if 15 != highestTreeIndex(8) {
		t.Errorf("highest tree index: %d", highestTreeIndex(8))
	}
}

func TestTreeHeadIndexes(t *testing.T) {
	idxs := treeHeadIndexes(11)
	var ids = []uint64{15, 18, 19}
	if !reflect.DeepEqual(idxs, ids) {
		t.Errorf("tree head indexes: %v", idxs)
	}

	idxs8 := treeHeadIndexes(8)
	var ids8 = []uint64{15}
	if !reflect.DeepEqual(idxs8, ids8) {
		t.Errorf("tree head indexes: %v", idxs8)
	}

	idxs5 := treeHeadIndexes(5)
	var ids5 = []uint64{7, 8}
	if !reflect.DeepEqual(idxs5, ids5) {
		t.Errorf("tree head indexes: %v", idxs5)
	}
	idxs0 := treeHeadIndexes(0)
	if len(idxs0) != 0 {
		t.Errorf("tree head indexes 0")
	}
}

func TestSubTreeSize(t *testing.T) {
	var s uint64 = 5
	ss := subTreeSize(s)
	if len(ss) != 2 || ss[0] != 7 || ss[1] != 1 {
		t.Errorf("subtree size: %d, %d", ss[0], ss[1])
	}
	var s0 uint64 = 0
	ss0 := subTreeSize(s0)
	if len(ss0) != 0 {
		t.Errorf("subtree size: %d", len(ss0))
	}
	var s8 uint64 = 8
	ss8 := subTreeSize(s8)
	if len(ss8) != 1 || ss8[0] != 15 {
		t.Errorf("subtree size: %d", ss0)
	}
}

func TestEmptyHash(t *testing.T) {
	var em = emptyHash()
	var emh = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if !bytes.Equal([]byte(fmt.Sprintf("%x", em.Bytes())), []byte(emh)) {
		t.Errorf("empty hash:\n%X", em.Bytes())
	}
}

func TestLeafHash(t *testing.T) {
	var l = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	tmp := append([]byte{0}, []byte(l)...)
	var h = leafHash([]byte(l))
	var s = sha256.Sum256(tmp)
	if !bytes.Equal(s[:], h.Bytes()) {
		t.Errorf("hash:\n%X\n%X", h.Bytes(), s[:])
	}
}

func TestNodeHash(t *testing.T) {
	var l = sha256.Sum256([]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"))
	var r = sha256.Sum256([]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b856"))
	var ln common.Uint256
	var rn common.Uint256
	if err := ln.FromBytes(l[:]); err != nil {
		t.Errorf("uint256: %s", err)
	}
	if err := rn.FromBytes(r[:]); err != nil {
		t.Errorf("uint256: %s", err)
	}
	rn.FromBytes(r[:])
	tmp := append(append([]byte{1}, l[:]...), r[:]...)
	var s = sha256.Sum256(tmp)
	var h = nodeHash(ln, rn)
	if !bytes.Equal(s[:], h.Bytes()) {
		t.Errorf("hash:\n%X\n%X", s[:], h[:])
	}
}

func TestReduceHash(t *testing.T) {
	var l1 = sha256.Sum256([]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"))
	var l2 = sha256.Sum256([]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b856"))
	var l3 = sha256.Sum256([]byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b857"))

	var u1 common.Uint256
	var u2 common.Uint256
	var u3 common.Uint256
	u1.FromBytes(l1[:])
	u2.FromBytes(l2[:])
	u3.FromBytes(l3[:])
	var ua = []common.Uint256{u1, u2, u3}
	var r = reduceHash(ua)

	l := len(ua)
	accum := ua[l-1]
	for i := l - 2; i >= 0; i-- {
		accum = nodeHash(ua[i], accum)
	}

	if !bytes.Equal(accum.Bytes(), r.Bytes()) {
		t.Errorf("reduce hash:\n%X\n%X", r.Bytes(), accum.Bytes())
	}

}
