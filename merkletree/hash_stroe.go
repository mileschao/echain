package merkletree

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/mileschao/echain/common"
)

var (
	// EmptyHash empty hash value
	EmptyHash = common.Uint256{}
	// ErrStoredHashLess the number of hash stored is less than expectation
	ErrStoredHashLess = errors.New("stored hashes are less than expected")
	// ErrStorageNil storage instance is nil
	ErrStorageNil = errors.New("storage is nil")
)

// HashStorage an interface for hash value storage
type HashStorage interface {
	Append(hash []common.Uint256) error
	Flush() error
	Close()
	GetHash(pos uint32) (common.Uint256, error)
}

// fileHashStorage an implementation of HashStorage interface
type fileHashStorage struct {
	fileName string
	file     *os.File
}

// NewFileHashStorage get HashStorage instance of file implement
func NewFileHashStorage(name string, leafSize uint64) (HashStorage, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	store := &fileHashStorage{
		fileName: name,
		file:     f,
	}

	err = store.checkConsistence(leafSize)
	if err != nil {
		return nil, err
	}

	num := totalStoredHashNum(leafSize)
	size := int64(num) * int64(common.UINT256_SIZE)

	_, err = store.file.Seek(size, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return store, nil
}

func totalStoredHashNum(leafSize uint64) uint64 {
	subtreesize := subTreeSize(leafSize)
	var sum uint64
	for _, v := range subtreesize {
		sum += uint64(v)
	}
	return sum
}

func (fhs *fileHashStorage) checkConsistence(leafSize uint64) error {
	num := totalStoredHashNum(leafSize)

	stat, err := fhs.file.Stat()
	if err != nil {
		return err
	} else if stat.Size() < int64(num)*int64(common.UINT256_SIZE) {
		return ErrStoredHashLess
	}
	return nil
}

// Append implement HashStorage interface
// append hash array into storage
func (fhs *fileHashStorage) Append(hash []common.Uint256) error {
	if fhs.file == nil { // do not store it
		return nil
	}
	b := new(bytes.Buffer)
	for _, h := range hash {
		h.Serialize(b)
	}
	_, err := fhs.file.Write(b.Bytes())
	return err
}

// Flush implement HashStorage interface
// flush file
func (fhs *fileHashStorage) Flush() error {
	if fhs.file == nil {
		return nil
	}
	return fhs.file.Sync()
}

// Close implement HashStorage interface
// Close file
func (fhs *fileHashStorage) Close() {
	if fhs.file == nil {
		return
	}
	fhs.file.Close()
}

// GetHash get hash value in storage by position
// `position` means the order of hash in storage
func (fhs *fileHashStorage) GetHash(pos uint32) (common.Uint256, error) {
	if fhs.file == nil {
		return EmptyHash, ErrStorageNil
	}
	hash := EmptyHash
	_, err := fhs.file.ReadAt(hash[:], int64(pos)*int64(common.UINT256_SIZE))
	if err != nil {
		return EmptyHash, err
	}
	return hash, nil
}

type memoryHashStorage struct {
	hashes []common.Uint256
}

func (mhs *memoryHashStorage) Append(hash []common.Uint256) error {
	mhs.hashes = append(mhs.hashes, hash...)
	return nil
}

func (mhs *memoryHashStorage) Flush() error {
	return nil
}

func (mhs *memoryHashStorage) Close() {
}

func (mhs *memoryHashStorage) GetHash(pos uint32) (common.Uint256, error) {
	if pos >= uint32(len(mhs.hashes)) {
		return common.UINT256_EMPTY, errors.New("memory hash storage out of range")
	}
	return mhs.hashes[pos], nil
}
