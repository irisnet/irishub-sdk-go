package sim

import (
	"github.com/irisnet/irishub-sdk-go/client"
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/distr"
	"github.com/irisnet/irishub-sdk-go/modules/event"
	"github.com/irisnet/irishub-sdk-go/modules/gov"
	"github.com/irisnet/irishub-sdk-go/modules/stake"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	addr    = "faa1d3mf696gvtwq2dfx03ghe64akf6t5vyz6pe3le"
	valAddr = "iva1x3f572u057lv88mva2q3z40ls8pup9hsg0lxcp"
	privKey = "927be78a5f5b63bb95ff34ed9c6e4b39b6af6d2f9f59731452de659cac9b19db"
)

type TestClient struct {
	bank.Bank
	event.Event
	stake.Stake
	gov.Gov
	distr.Distr
}

func NewTestClient() TestClient {
	cdc := types.NewAmino()
	rpc := net.NewRPCClient("localhost:26657")
	ctx := &types.TxContext{
		Codec:   cdc,
		ChainID: "irishub-test",
		Online:  true,
		KeyDAO:  createTestKeyDAO(),
		Network: types.Testnet,
		Mode:    types.Commit,
		RPC:     rpc,
	}
	txm := client.NewBaseClient(ctx)
	return TestClient{
		Bank:  bank.New(txm),
		Stake: stake.New(txm),
		Event: event.New(txm),
		Gov:   gov.New(txm),
		Distr: distr.New(txm),
	}
}

func (TestClient) GetTestSender() string {
	return addr
}

func (TestClient) GetTestValidator() string {
	return valAddr
}

func createTestKeyDAO() TestKeyDAO {
	dao := TestKeyDAO{
		store: map[string]types.KeyStore{},
	}
	keystore := TestKeystore{
		Private: privKey,
		Address: addr,
	}
	_ = dao.Write("test1", keystore)
	return dao
}

type TestKeyDAO struct {
	store map[string]types.KeyStore
}

func (dao TestKeyDAO) Write(name string, keystore types.KeyStore) error {
	dao.store[name] = keystore
	return nil
}

func (dao TestKeyDAO) Read(name string) types.KeyStore {
	return dao.store[name]
}

func (dao TestKeyDAO) Delete(name string) error {
	return nil
}

type TestKeystore struct {
	Private string
	Address string
}

func (t TestKeystore) GetPrivate() string {
	return t.Private
}
func (t TestKeystore) GetAddress() string {
	return t.Address
}
