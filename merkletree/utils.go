package merkletree

import (
	"crypto/sha256"

	"github.com/mileschao/echain/common"
)

var (
	// HashLeafPrefix prefix to hash leaf byte array
	HashLeafPrefix = []byte{0}
	// HashNodePrefix prefix to hash left and right children node
	HashNodePrefix = []byte{1}
)

func countBit(num uint64) uint32 {
	var count uint32
	for num != 0 {
		num &= (num - 1)
		count++
	}
	return count
}

func isPower2(num uint64) bool {
	return countBit(num) == 1
}

// highest bit 1 index
// i.e 0000 0000 0000 0000 0000 0000 0001 0000
//                                      ^
// highBit = 5
func highBit(num uint64) uint64 {
	var hiBit uint64
	for num != 0 {
		num >>= 1
		hiBit++
	}
	return hiBit
}

func lowBit(num uint64) uint64 {
	return highBit(num & -num)
}

func isEvenNumber(num uint64) bool {
	return num%2 == 0
}

func isOddNumber(num uint64) bool {
	return !isEvenNumber(num)
}

// calculate merkle tree's height by leaf size
// i.e
//         H              -->height = 3
//     /        \
//    H         H
//  /   \     /   \
// 1    2    3    4    5  -->leafSize = 4
func forestHeight(leafSize uint64) uint64 {
	return highBit(leafSize)
}

// the Merkle tree can be seen as an `forest` that contains many `tree`
// the highest tree is always the leftest side tree
// and the tree number in the `forest` is the bit 1 count number of leafSize
// i.e. the Merkle tree with leafSize = 5 = 0000 0101
// the count bit is 2, and the number of tree in forest is also 2
func treeCountInForest(leafSize uint64) uint32 {
	return countBit(leafSize)
}

// calculate the highest tree's head index (base on 1) according to the hight that deduced by leafSize
// that is:
// 			leafSize --> height --> highest tree head index
// i.e
//         H           --> highest tree head
//     /       \
//    H         H
//  /   \     /   \
// 1    2    3    4   5
// ^^^^^^^^^^^^^^^^
// highest tree
//
// the highest tree is always the leftest side tree in the forest
// thus we can get the highest tree head index by height
// i.e. height = 3 means the number of node in every level is:
// level 3: 0000 0001 = 1
// level 2: 0000 0010 = 2
// level 1: 0000 0100 = 4 = (0x01 << height) is the leafsize
func highestTreeIndex(leafSize uint64) uint64 {
	h := forestHeight(leafSize)
	var index uint64
	for i := uint64(0); i < h; i++ {
		index += (0x01 << i)
	}
	return index
}

// get index array that contains every tree head index(base on 1) in forest
//
//
//                 H(index 15)
//           /                 \
//         H                    H
//     /        \          /        \
//    H         H         H         H         H (index 18)
//  /   \     /   \     /   \     /   \     /   \
// 1    2    3    4    5    6    7    8    9    10    11
// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^    ^^^^^^^    ^^
// treeCount = 3
// preTreeLeafSize      --> 8 --> 10 --> 11
// leafSizeleft         --> 3 --> 1  --> 0
// currentTreeLeafSize  --> 8 --> 2  --> 1
// indexes              -->[15] -->[15, 18] --> [15, 18, 19]
func treeHeadIndexes(leafSize uint64) []uint64 {
	treeSizes := subTreeSize(leafSize)
	var treeIndexes = make([]uint64, 0, len(treeSizes))
	treeIndexes = append([]uint64{0}, treeIndexes...) // trick to do reduce
	for i := 0; i < len(treeSizes); i++ {
		treeIndexes = append(treeIndexes, treeIndexes[i]+treeSizes[i])
	}
	return treeIndexes[1:]
}

// subTreeSize get sub tree size
// i.e.
//           H
//       /      \
//      H       H
//    /  \    /  |
//   H   H   H   H   H   ---> leaf nodes
//  ^^^^^^^^^^^^^^   ^
//   subtree_1        subtree_2
// leafSize = 5, countBit(leafSize) is the count number of subtree
// subTreeSize(leafSize) = [7, 1]
func subTreeSize(leafSize uint64) []uint64 {
	treeCount := treeCountInForest(leafSize)
	indexes := make([]uint64, 0, treeCount)
	preTreeLeafSize := uint64(0)
	for i := uint32(0); i < treeCount; i++ {
		leafSizeleft := leafSize - preTreeLeafSize
		currentTreeLeafSize := uint64(0x01 << (forestHeight(leafSizeleft) - 1))
		indexes = append(indexes, highestTreeIndex(currentTreeLeafSize))
		preTreeLeafSize += currentTreeLeafSize
	}
	return indexes
}

// emptyHash return fix result:
// e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
func emptyHash() common.Uint256 {
	var u common.Uint256
	var e = sha256.Sum256(nil)
	u.FromBytes(e[:])
	return u
}

func leafHash(l []byte) common.Uint256 {
	// TODO: the param type []byte? or common.Uint256?
	tmp := append(HashLeafPrefix, l...)
	b := sha256.Sum256(tmp)
	var u common.Uint256
	u.FromBytes(b[:])
	return u
}

// nodeHash sha256 with left child node and right child node
func nodeHash(left common.Uint256, right common.Uint256) common.Uint256 {
	data := append(append(HashNodePrefix, left[:]...), right[:]...)
	b := sha256.Sum256(data)
	var u common.Uint256
	u.FromBytes(b[:])
	return u
}

// reduceHash do reduce `nodeHash` on the hash array
func reduceHash(h []common.Uint256) common.Uint256 {
	return reduce(h[:], nodeHash)
}

// reduce auxiliary function partialy implement `reduce` in functional programing
// i.e
// reduce([1,2,3,4,5,6]) with reducer `add` function
// means 6 + 5 + 4 + 3 + 2 + 1
func reduce(arr []common.Uint256, reducer func(common.Uint256, common.Uint256) common.Uint256) common.Uint256 {
	if len(arr) == 0 {
		return emptyHash()
	} else if len(arr) == 1 {
		return arr[0]
	} else { // >= 1
		n := arr[len(arr)-1]
		for i := len(arr) - 2; i >= 0; i-- {
			n = reducer(arr[i], n)
		}
		return n
	}
}
