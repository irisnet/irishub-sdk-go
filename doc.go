// Package client is the entrance of the entire SDK function. SDKConfig is used to configure SDK parameters.
//
// The SDK mainly provides the functions of the following modules, including:
// asset, bank, distribution, gov, keys, oracle, random, service, slashing, staking
//
// As a quick start:
//
// 	fees, err := types.ParseCoins("1iris")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	client := sdk.NewClient(types.SDKConfig{
// 		NodeURI:   NodeURI,
// 		Network:   Network,
// 		ChainID:   ChainID,
// 		Gas:       Gas,
// 		Fee:       fees,
// 		KeyDAO:    createTestKeyDAO(),
// 		Mode:      Mode,
// 		Online:    Online,
// 		StoreType: types.Key,
// 		Level:     "debug",
// 	})
// KeyDAO is an interface, you need to implement this interface to support the ability of SDK to access external data, such as access to a database
// 	func createTestKeyDAO() types.KeyDAO {
// 		dao := TestKeyDAO{
// 		     store: map[string]types.Store{},
// 		}
// 		keystore := types.KeyInfo{
// 		 		PrivKey: priKey,
// 		 		Address: addr,
// 		}
// 		_ = dao.Write("test1", keystore)
// 		return types.NewKeyDAO(&dao,nil)
//	}
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
package sdk
