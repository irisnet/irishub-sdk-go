// Package bank is mainly used to transfer coins between accounts,query account balances, and implement interface rpc.Bank
//
// In addition, the available units of tokens in the IRIShub system are defined using [coin-type](https://www.irisnet.org/docs/concepts/coin-type.html).
//
// [More Details](https://www.irisnet.org/docs/features/bank.html)
//
// As a quick start:
//
// Transfer coins to other account
//
//  client := test.NewClient()
//  amt := types.NewIntWithDecimal(1, 18)
//  coins := types.NewCoins(types.NewCoin("iris-atto", amt))
//  to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
//  baseTx := types.BaseTx{
// 		From: "test1",
// 		Gas:  20000,
// 		Memo: "test",
// 		Mode: types.Commit,
//  }
//  result,err := client.Bank().Send(to,coins,baseTx)
//  fmt.Println(result)
//
// Burn some coins from your account
//
//  client := sim.NewClient()
//  amt := types.NewIntWithDecimal(1, 18)
//  coins := types.NewCoins(types.NewCoin("iris-atto", amt))
//  baseTx := types.BaseTx{
// 		From: "test1",
// 		Gas:  20000,
// 		Memo: "test",
// 		Mode: types.Commit,
//  }
//  result,err := client.Bank().Burn(coins, baseTx)
//  fmt.Println(result)
//
// Set account memo
//
//  client := sim.NewClient()
//  result,err := client.Bank().SetMemoRegexp("testMemo", baseTx)
//  fmt.Println(result)
//
// Query account information
//
//  client := sim.NewClient()
//  result,err := client.Bank().QueryAccount("faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm")
//  fmt.Println(result)
//
// Query the token information
//
//  client := sim.NewClient()
//  result,err := client.Bank().QueryTokenStats("iris")
//  fmt.Println(result)
//
package bank
