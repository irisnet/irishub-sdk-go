package types

import (
	"errors"
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

// NewKeyBase initialize a keybase based on the configuration
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
	existed, _ := k.Has(name)
	if existed {
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
	if err != nil || bz == nil {
		return store, errors.New(fmt.Sprintf("key %s not exist", name))
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
func (k KeyBase) Has(name string) (bool, error) {
	return k.db.Has(infoKey(name))
}

func infoKey(name string) []byte {
	return []byte(fmt.Sprintf("%s.%s", name, infoSuffix))
}
