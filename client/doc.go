// Package client is the entrance of the entire SDK function. SDKConfig is used to configure SDK parameters.
//
// The SDK mainly provides the functions of the following modules, including:
// asset, bank, distribution, gov, keys, oracle, random, service, slashing, staking
//
// As a quick start:
//
// 	keyManager, err := crypto.NewMnemonicKeyManager(mnemonic)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	priKey, err = keyManager.ExportAsPrivateKey()
// 	if err != nil {
// 		panic(err)
// 	}
// 	addr = types.AccAddress(keyManager.GetPrivKey().PubKey().Address()).String()
// 	fees, err := types.ParseCoins(Fee)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	client := client.NewSDKClient(types.SDKConfig{
// 		NodeURI:   NodeURI,
// 		Network:   Network,
// 		ChainID:   ChainID,
// 		Gas:       Gas,
// 		Fee:       fees,
// 		KeyDAO:    createTestKeyDAO(),
// 		Mode:      Mode,
// 		Online:    Online,
// 		StoreType: types.Keystore,
// 		Level:     "debug",
// 	})
// KeyDAO is an interface, you need to implement this interface to support the ability of SDK to access external data, such as access to a database
// 	func createTestKeyDAO() TestKeyDAO {
// 		dao := TestKeyDAO{
// 		store: map[string]types.Store{},
// 	}
// 	keystore := types.KeyInfo{
// 		PrivKey: priKey,
// 		Address: addr,
// 	}
// 	_ = dao.Write("test1", keystore)
// 		return dao
// 	}
//
// 	type TestKeyDAO struct {
// 		store map[string]types.Store
// 	}
//
// 	func (dao TestKeyDAO) Write(name string, store types.Store) error {
// 		dao.store[name] = store
// 		return nil
// 	}
//
// 	func (dao TestKeyDAO) Read(name string) (types.Store, error) {
// 		return dao.store[name], nil
// 	}
//
// 	func (dao TestKeyDAO) Delete(name string) error {
// 		return nil
// 	}
package client
