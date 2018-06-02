package merkletree

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/mileschao/echain/common"
)

var (
	// ErrBadLeafSize the leafsize is not match with hash list
	ErrBadLeafSize = errors.New("number of hashes do not match number of bit in leaf size")
)

// MerkleHeap Merkle Tree's corresponding heap format
// just as the 'binary heap' to the 'binary tree',
// MerkleHeap is an array that easy to serialize 'Merkle Tree' into buffer
// i.e.
// merkle Tree:
//
//        H(H(H12,H34),H(H12,H34))^
//           /              \
//    * H(H12,H34)      * H(H56,H78)
//     /       \         /      \
// * H12     * H34    * H56   * H78
//  /   \     /   \    /  \    /  \
//*1   *2   *3   *4  *5  *6  *7  *8   9^
//
// corresponding merkle heap:
// 1|2|H12|3|4|H34|H(H12,H34)|5|6|H56|7|8|H78|H(H56,H78)|H(H(H12,H34),H(H12,H34))|9
// MerkleHeap.heigth = 4
// MerkleHeap.upperNodes = [H(H(H12,H34),H(H12,H34)), 9]
// MerkleHeap.root = H(H(H12,H34),H(H12,H34))
// MerkleHeap.leafSize = 9
//
// there are some trick here:
// 1. when storing into database the merkle heap above is divided into two parts:
//
// part1: elements with `*` are stored into HashStorage(i.e fileHashStorage)
// these elements with `*` means that all those has already been hashed
// these elements in above example are: 1|2|H12|3|4|H34|H(H12,H34)|5|6|H56|7|8|H78|H(H56,H78)
//
// part2: elements with `^` are stored into stateStore in ledger(i.e leveldb)
// these elements with `^` means that all those has not been hashed
// these elements in above example are: H(H(H12,H34),H(H12,H34))|9
// these are always the upper node of the merkle tree
//
// 2. tricks on MerkleHeap.leafSize
// 2.1 the leafSize means the number of leaves in Merkle Tree
// 		i.e the leafSize in the example above is 9
// 2.2 the bit 1 count of leafSize is the number of elements with `^`
//		i.e when leafSize = 9 = 0000 1001; countBit(leafSize) = 2
//			the elements with `^` in above example is: H(H(H12,H34),H(H12,H34))|9
// 2.3 the highest bit 1 index of leafSize is the height of the Merkle Tree
//		i.e when leafSize = 9 == 0000 1001; hightBit(leafSize) = 4
//			the Merkle Tree's height is 4
// 3. MerkleHeap.upperNodes is the list of element that to be hashed
//		i.e the upperNodes is the elements with `^` in above example: H(H(H12,H34),H(H12,H34))|9
//			and is going to stored into stateStore of ledger in futher
type MerkleHeap struct {
	height      uint64
	upperNodes  []common.Uint256
	hashStorage HashStorage
	root        common.Uint256
	leafSize    uint64
}

// NewMerkleStorage create an new merkle storage with MerkleHeap data struct
// BUG(r): if NewMerkleStorage with leafSize != 0, the AddLeaf will lose the data
func NewMerkleStorage(leafSize uint64, hashes []common.Uint256, store HashStorage) *MerkleHeap {

	mh := &MerkleHeap{
		height:      0,
		upperNodes:  nil,
		hashStorage: store,
		root:        EmptyHash,
	}

	mh.update(leafSize, hashes)
	return mh
}

// UpperNodes get the upper node list of merkle heap that to be hashed
func (mh *MerkleHeap) UpperNodes() []common.Uint256 {
	return mh.upperNodes
}

// LeafSize get the leaf number of the merkle tree
func (mh *MerkleHeap) LeafSize() uint64 {
	return mh.leafSize
}

// Root get the root hash of the merkle tree
func (mh *MerkleHeap) Root() common.Uint256 {
	return mh.root
}

// RootWithNewLeaf calculate hash of merkle tree's upper node list and leaf input
func (mh *MerkleHeap) RootWithNewLeaf(leaf common.Uint256) common.Uint256 {
	return reduceHash(append(mh.upperNodes, leaf))
}

// AddLeaf add new leaf node hash to Merkle Tree
// when an new leaf add into merkle tree
// new hash level up maybe calculate
// thus, the hased nodes need to be stored into the hashStorage
// and, the upper node list need to be update
// i.e
//                                                          countBit(leafSize) = len(upperNodes) = 2
//                                                         /hightBit(leafSize) = height = 2
//                                                        v
//    H12 --->uppernode         |   leafSize = 3 = 0000 0011
//   /   |                      |   upperNodes = [H12, 3]
//  1    2     3 -->uppernode   |   height = 2
//  ^    ^                      |
//  |----|--->hashStorage       |
// ========================================================
// when Add new leaf `4` into MerkleTree:
//                                                            countBit(leafSize) = len(uppernode) = 1
//                                                           /hightBit(leafSize) = height = 3
//                                                          v
//      H(H12,H34)   ---> uppernode  | leafSize = 4 = 0000 0100
//      /       \                    | upperNodes = [H(H12, H34)]
//    H12       H34  --> hashStorage | height = 3
//   /   \     /   \                 |
//  1    2    3    4 --> hashStorage |
//
//
func (mh *MerkleHeap) AddLeaf(leaf common.Uint256) []common.Uint256 {
	// TODO: ugly code, make it better to read for human beings
	size := len(mh.upperNodes)
	auditPath := make([]common.Uint256, size, size)
	nodeToStored := make([]common.Uint256, 0)
	// reverse
	for i, v := range mh.upperNodes {
		auditPath[size-i-1] = v
	}
	// BUG(r): if the pre leafSize if even number and the upperNodes is not empty
	// the data in upperNode will lose, as to the BUG in `NewMerkleStorage`
	nodeToStored = append(nodeToStored, leaf)
	//mh.height = 1
	for s := mh.leafSize; s%2 == 1; s = s >> 1 { // odd number
		//mh.height++
		leaf = nodeHash(mh.upperNodes[size-1], leaf)
		nodeToStored = append(nodeToStored, leaf)
		size--
	}
	if mh.hashStorage != nil {
		mh.hashStorage.Append(nodeToStored)
		mh.hashStorage.Flush()
	}
	mh.leafSize++
	mh.upperNodes = mh.upperNodes[0:size]
	mh.upperNodes = append(mh.upperNodes, leaf)
	mh.root = EmptyHash
	mh.height = highBit(mh.leafSize)

	return auditPath
}

// Proofs get proofs
// m is zero base leaf index, n is leafsize(thus, 1 based)
// what is proof?
// i.e. the merkle tree below proof is the node with `*`
//
//                    H
//           /                 \
//          H*                 H
//     /        \          /        \
//    H         H         H         H*        H ------------->Hash(H, 11)*
//  /   |     /   \     /   \     /   \     /   \          |
// 1    2    3    4    5    6*   7    8    9    10    11---|
//                     ^                               ^
//                     m = index = 4                   n = leafSize = 11
//
// when m = 4 && n = 11
// proofs = [6, H(7, 8), H(H12, H34), H(H910, 11)]
// the principle is:
// 1. the tree of forest that m belong to, divide the tree into server sub trees with out node m,
//		the proof is the root node of these subtrees order by their height
//		in the example above, they are: 6, H(7, 8), H(H12, H34)
// 2. the right side trees of forest that the tree m belong to, do hash of all the root node of these trees
//		in the example above is H(H910, 11)
func (mh *MerkleHeap) Proofs(m, n uint64) ([]common.Uint256, error) {
	//TODO: ugly code, make it better to read for human being
	if m >= n {
		return nil, errors.New("wrong parameters")
	} else if mh.leafSize < n {
		return nil, errors.New("not available yet")
	} else if mh.hashStorage == nil {
		return nil, errors.New("hash store not available")
	}

	offset := uint64(0)
	var hashes []common.Uint256
	for n != 1 {
		k := uint64(1 << (highBit(n-1) - 1))
		if m < k {
			pos := treeHeadIndexes(n - k)
			subhashes := make([]common.Uint256, len(pos), len(pos))
			for p := range pos {
				pos[p] += offset + k*2 - 1
				subhashes[p], _ = mh.hashStorage.GetHash(uint32(pos[p] - 1))
			}
			rootk2n := reduceHash(subhashes)
			hashes = append(hashes, rootk2n)
			n = k
		} else {
			offset += k*2 - 1
			root02k, _ := mh.hashStorage.GetHash(uint32(offset - 1))
			hashes = append(hashes, root02k)
			m -= k
			n -= k
		}
	}

	length := len(hashes)
	reverse := make([]common.Uint256, length, length)
	for k := range reverse {
		reverse[k] = hashes[length-k-1]
	}

	return reverse, nil
}

// Serialize implement the common.Serialzable interface
func (mh *MerkleHeap) Serialize(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, mh.leafSize); err != nil {
		return err
	}
	for _, h := range mh.upperNodes {
		if err := h.Serialize(w); err != nil {
			return err
		}
	}
	return nil
}

// Deserialize implement the common.Serializable interface
func (mh *MerkleHeap) Deserialize(r io.Reader) error {
	var leafSize uint64
	if err := binary.Read(r, binary.LittleEndian, &leafSize); err != nil {
		return err
	}
	num := countBit(leafSize)
	var upperNodes []common.Uint256
	for i := uint32(0); i < num; i++ {
		if err := upperNodes[i].Deserialize(r); err != nil {
			return err
		}
	}
	return nil
}

// update merkle heap with leaf size and uppernodes list
func (mh *MerkleHeap) update(leafSize uint64, toBeHashed []common.Uint256) error {
	bitCount := countBit(leafSize)
	if len(toBeHashed) != int(bitCount) {
		panic("number of hashes != num bit in tree_size")
	}
	mh.leafSize = leafSize
	mh.upperNodes = toBeHashed
	mh.height = highBit(leafSize)
	mh.root = EmptyHash
	return nil
}
