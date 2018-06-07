package leveldb

import (
	"github.com/mileschao/echain/storage"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//Iterator leveldb iterator
type Iterator struct {
	iter iterator.Iterator
}

//Next implement storage.Iterator implement
func (it *Iterator) Next() bool {
	return it.iter.Next()
}

//Prev implement storage.Iterator implement
func (it *Iterator) Prev() bool {
	return it.iter.Prev()
}

//First implement storage.Iterator implement
func (it *Iterator) First() bool {
	return it.iter.First()
}

//Last implement storage.Iterator implement
func (it *Iterator) Last() bool {
	return it.iter.Last()
}

//Seek implement storage.Iterator implement
func (it *Iterator) Seek(key []byte) bool {
	return it.iter.Seek(key)
}

//Key implement storage.Iterator implement
func (it *Iterator) Key() []byte {
	return it.iter.Key()
}

//Value implement storage.Iterator implement
func (it *Iterator) Value() []byte {
	return it.iter.Value()
}

//Release implement storage.Iterator implement
func (it *Iterator) Release() {
	it.iter.Release()
}

//Storage leveldb storage
// implement PersistStorage interface
type Storage struct {
	db    *leveldb.DB // LevelDB instance
	batch *leveldb.Batch
}

// used to compute the size of bloom filter bits array .
// too small will lead to high false positive rate.
const BITSPERKEY = 10

//NewStore return LevelDBStore instance
func NewStore(file string) (*Storage, error) {

	// default Options
	o := opt.Options{
		NoSync: false,
		Filter: filter.NewBloomFilter(BITSPERKEY),
	}

	db, err := leveldb.OpenFile(file, &o)

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}

	if err != nil {
		return nil, err
	}

	return &Storage{
		db:    db,
		batch: nil,
	}, nil
}

//Put implement Persist storage interface
func (s *Storage) Put(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

//Get implement Persist storage interface
func (s *Storage) Get(key []byte) ([]byte, error) {
	dat, err := s.db.Get(key, nil)
	return dat, err
}

//Has implement Persist storage interface
func (s *Storage) Has(key []byte) (bool, error) {
	return s.db.Has(key, nil)
}

//Delete implement Persist storage interface
func (s *Storage) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

//NewBatch implement Persist storage interface
func (s *Storage) NewBatch() {
	s.batch = new(leveldb.Batch)
}

//BatchPut implement Persist storage interface
func (s *Storage) BatchPut(key []byte, value []byte) {
	s.batch.Put(key, value)
}

//BatchDelete implement Persist storage interface
func (s *Storage) BatchDelete(key []byte) {
	s.batch.Delete(key)
}

//BatchCommit implement Persist storage interface
func (s *Storage) BatchCommit() error {
	err := s.db.Write(s.batch, nil)
	if err != nil {
		return err
	}
	s.batch = nil
	return nil
}

//Close implement Persist storage interface
func (s *Storage) Close() error {
	err := s.db.Close()
	return err
}

//NewIterator implement Persist storage interface
func (s *Storage) NewIterator(prefix []byte) storage.Iterator {

	iter := s.db.NewIterator(util.BytesPrefix(prefix), nil)

	return &Iterator{
		iter: iter,
	}
}
