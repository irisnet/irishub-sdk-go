package sim

import (
	"github.com/irisnet/irishub-sdk-go/client"
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/event"
	"github.com/irisnet/irishub-sdk-go/modules/stake"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	addr    = "iaa1x3f572u057lv88mva2q3z40ls8pup9hsa74f9x"
	valAddr = "iva1x3f572u057lv88mva2q3z40ls8pup9hsg0lxcp"
	privKey = "8D03FEDB094224959DD12016D24782429216246BC03084211C0305F9767C3C38"
)

type TestClient struct {
	Bank  bank.Client
	Event event.Event
	Stake stake.Client
}

func NewTestClient() TestClient {
	cdc := types.NewAmino()
	rpc := net.NewRPCClient("localhost:26657")
	ctx := &types.TxContext{
		Codec:   cdc,
		ChainID: "irishub-test",
		Online:  true,
		KeyDAO:  createTestKeyDAO(),
		Network: types.Mainnet,
		Mode:    types.Commit,
		RPC:     rpc,
	}
	txm := client.NewBaseClient(ctx)
	return TestClient{
		Bank:  bank.NewClient(txm),
		Stake: stake.NewClient(txm),
		Event: event.NewEvent(txm),
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
