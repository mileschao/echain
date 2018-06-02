package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/mileschao/echain/common"
)

var (
	//ErrMerkleTreeEmpty error that try to construct merkle tree with 0 item
	ErrMerkleTreeEmpty = errors.New("construct merkle tree with 0 item")
)

type merkleTreeNode struct {
	Hash  common.Uint256
	Left  *merkleTreeNode
	Right *merkleTreeNode
}

type merkleTree struct {
	Depth uint64
	Root  *merkleTreeNode
}

func repeatSha256(s []common.Uint256) common.Uint256 {
	b := new(bytes.Buffer)
	for _, d := range s {
		d.Serialize(b)
	}
	temp := sha256.Sum256(b.Bytes())
	f := sha256.Sum256(temp[:])
	return common.Uint256(f)
}

//generate the up level merkle tree nodes with repeat sha256
func levelUp(nodes []*merkleTreeNode) []*merkleTreeNode {
	var nextLevel []*merkleTreeNode
	for i := 0; i < len(nodes)/2; i++ {
		var data []common.Uint256
		data = append(data, nodes[i*2].Hash)   // even number index
		data = append(data, nodes[i*2+1].Hash) // odd number index
		hash := repeatSha256(data)
		node := &merkleTreeNode{
			Hash:  hash,
			Left:  nodes[i*2],
			Right: nodes[i*2+1],
		}
		nextLevel = append(nextLevel, node)
	}
	if len(nodes)%2 == 1 { // odd number: the last one repeat sha256 with itself
		var data []common.Uint256
		data = append(data, nodes[len(nodes)-1].Hash)
		data = append(data, nodes[len(nodes)-1].Hash)
		hash := repeatSha256(data)
		node := &merkleTreeNode{
			Hash:  hash,
			Left:  nodes[len(nodes)-1],
			Right: nodes[len(nodes)-1],
		}
		nextLevel = append(nextLevel, node)
	}
	return nextLevel
}

func newLeaves(hashes []common.Uint256) []*merkleTreeNode {
	var leaves []*merkleTreeNode
	for _, h := range hashes {
		node := &merkleTreeNode{
			Hash: h,
		}
		leaves = append(leaves, node)
	}
	return leaves
}

//create a new merkleTree by hash array
// show as below:
//         root        |         <- level up
//       /     \       |
//      n      n       |         <- level up
//    /    \   ||      | heigth
//   n     n    n      |         <- level up
//  / \   / \   ||     |
// l  l  l  l    l     |         <- input 5 hash array
func newMerkleTree(hashes []common.Uint256) (*merkleTree, error) {
	if len(hashes) == 0 {
		return nil, ErrMerkleTreeEmpty
	}
	var height uint64

	height = 1
	nodes := newLeaves(hashes)
	for len(nodes) > 1 {
		nodes = levelUp(nodes)
		height++
	}
	mt := &merkleTree{
		Root:  nodes[0],
		Depth: height,
	}
	return mt, nil

}

// CalcMerkleTreeRoot calculate merkle tree root hash by hash array
func CalcMerkleTreeRoot(hashes []common.Uint256) common.Uint256 {
	if len(hashes) == 0 {
		return common.UINT256_EMPTY
	}
	if len(hashes) == 1 {
		return hashes[0]
	}
	tree, _ := newMerkleTree(hashes)
	return tree.Root.Hash
}
