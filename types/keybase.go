package types

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
	_ KeyDAO = KeyBase{}
)

type KeyBase struct {
	db  dbm.DB
	cdc Codec
	AES
}

// NewKeyBase initialize a keybase based on the configuration.
// Use leveldb as storage
func NewKeyBase(rootDir string, cdc Codec) (KeyDAO, error) {
	db, err := dbm.NewGoLevelDB(keyDBName, filepath.Join(rootDir, "keys"))
	if err != nil {
		return nil, err
	}
	keybase := KeyBase{
		db:  db,
		cdc: cdc,
	}
	return keybase, nil
}

// Write add a key information to the local store
func (k KeyBase) Write(name string, store Store) error {
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
func (k KeyBase) Read(name string) (store Store, err error) {
	bz, err := k.db.Get(infoKey(name))
	if bz == nil || err != nil {
		return store, err
	}

	if err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &store); err != nil {
		return store, err
	}
	return
}

// Delete delete a key from the local store
func (k KeyBase) Delete(name string) error {
	return k.db.DeleteSync(infoKey(name))
}

// Delete delete a key from the local store
func (k KeyBase) Has(name string) bool {
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
type Memory struct {
	store map[string]Store
	AES
}

func NewMemory() Memory {
	return Memory{
		store: make(map[string]Store),
	}
}
func (m Memory) Write(name string, store Store) error {
	m.store[name] = store
	return nil
}

func (m Memory) Read(name string) (Store, error) {
	return m.store[name], nil
}

func (m Memory) Delete(name string) error {
	delete(m.store, name)
	return nil
}

func (m Memory) Has(name string) bool {
	_, ok := m.store[name]
	return ok
}
