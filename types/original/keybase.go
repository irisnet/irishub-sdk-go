package original

import (
	"fmt"
	"path/filepath"

	dbm "github.com/tendermint/tm-db"
)

const (
	keyDBName  = "keys"
	infoSuffix = "info"
)

var (
	_ KeyDAO = LevelDB{}
)

type LevelDB struct {
	db  dbm.DB
	cdc Codec
	AES
}

// NewLevelDB initialize a keybase based on the configuration.
// Use leveldb as storage
func NewLevelDB(rootDir string, cdc Codec) (KeyDAO, error) {
	db, err := dbm.NewGoLevelDB(keyDBName, filepath.Join(rootDir, "keys"))
	if err != nil {
		return nil, err
	}
	keybase := LevelDB{
		db:  db,
		cdc: cdc,
	}
	return keybase, nil
}

// Write add a key information to the local store
func (k LevelDB) Write(name string, store Store) error {
	if k.Has(name) {
		return fmt.Errorf("name %s has exist", name)
	}

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(store)
	if err != nil {
		return err
	}
	return k.db.SetSync(infoKey(name), bz)
}

// Read read a key information from the local store
func (k LevelDB) Read(name string) (store Store, err error) {
	bz, err := k.db.Get(infoKey(name))
	if bz == nil || err != nil {
		return store, err
	}

	if err := k.cdc.UnmarshalBinaryBare(bz, &store); err != nil {
		return store, err
	}
	return
}

// Delete delete a key from the local store
func (k LevelDB) Delete(name string) error {
	return k.db.DeleteSync(infoKey(name))
}

// Delete delete a key from the local store
func (k LevelDB) Has(name string) bool {
	existed, err := k.db.Has(infoKey(name))
	if err != nil {
		return false
	}
	return existed
}

func infoKey(name string) []byte {
	return []byte(fmt.Sprintf("%s.%s", name, infoSuffix))
}

// Use memory as storage, use with caution in build environment
type MemoryDB struct {
	store map[string]Store
	AES
}

func NewMemoryDB() MemoryDB {
	return MemoryDB{
		store: make(map[string]Store),
	}
}
func (m MemoryDB) Write(name string, store Store) error {
	m.store[name] = store
	return nil
}

func (m MemoryDB) Read(name string) (Store, error) {
	return m.store[name], nil
}

func (m MemoryDB) Delete(name string) error {
	delete(m.store, name)
	return nil
}

func (m MemoryDB) Has(name string) bool {
	_, ok := m.store[name]
	return ok
}
