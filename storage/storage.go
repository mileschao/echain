package storage

import (
	"github.com/mileschao/echain/common/serialize"
)

type ItemState byte

//Status of item
const (
	None    ItemState = iota //no change
	Changed                  //which was be mark delete
	Deleted                  //which wad be mark delete
)

//State item struct
type StateItem struct {
	Key   string                 //State key
	Value serialize.Serializable //State value
	State ItemState              //Status
	Trie  bool                   //no use
}

func (e *StateItem) Copy() *StateItem {
	c := *e
	return &c
}

//PersistStorage persistent storage
type PersistStorage interface {
	Put(key []byte, value []byte) error //Put the key-value pair to store
	Get(key []byte) ([]byte, error)     //Get the value if key in store
	Has(key []byte) (bool, error)       //Whether the key is exist in store
	Delete(key []byte) error            //Delete the key in store
	NewBatch()                          //Start commit batch
	BatchPut(key []byte, value []byte)  //Put a key-value pair to batch
	BatchDelete(key []byte)             //Delete the key in batch
	BatchCommit() error                 //Commit batch to store
	Close() error                       //Close store
	NewIterator(prefix []byte) Iterator //Return the iterator of store
}

//Iterator iterator of storage
type Iterator interface {
	Next() bool           //Next item. If item available return true, otherwise return false
	Prev() bool           //previous item. If item available return true, otherwise return false
	First() bool          //First item. If item available return true, otherwise return false
	Last() bool           //Last item. If item available return true, otherwise return false
	Seek(key []byte) bool //Seek key. If item available return true, otherwise return false
	Key() []byte          //Return the current item key
	Value() []byte        //Return the current item value
	Release()             //Close iterator
}

//MemoryStorage storage that in memory
type MemoryStorage interface {
	Put(prefix byte, key []byte, value serialize.Serializable, state ItemState) //Put the key-value pair to store
	Get(prefix byte, key []byte) *StateItem                                     //Get the value if key in store
	Delete(prefix byte, key []byte)                                             //Delete the key in store
	GetChangeSet() map[string]*StateItem                                        //Get all updated key-value set
	Find() []*StateItem                                                         // Get all key-value in store
}
