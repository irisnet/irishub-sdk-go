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
// 	client := sdk.NewClient(types.ClientConfig{
// 		NodeURI:   NodeURI,
// 		Network:   Network,
// 		ChainID:   ChainID,
// 		Gas:       Gas,
// 		Fee:       fees,
// 		Mode:      Mode,
// 		Online:    Online,
// 		StoreType: types.PrivKey,
// 		Level:     "debug",
// 	})
package sdk
