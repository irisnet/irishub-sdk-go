# IRISHUB Chain Go SDK

Irishub Chain GO SDK 对Irishub提供的API做了一层简单的封装，为用户快速开发基于irishub chain的应用提供的极大的方便，它包括以下核心组件：

- adapter - 封装了和用户DAO层CRUD操作的基本操作
- crypto - 实现了irishub私钥、地址、keystore生成方法
- modules - 实现了rpc包下API接口，封装了irishub各模块的API方法。
- rpc - SDK向用户暴露的API接口
- test - 自动化测试入口
- types - 核心类型
- utils - 工具库的包装

其中`modules`包下，实现了irishub各模块对外提供的部分或者全部接口，目前主要包括：`asset`、`bank`、`distribution`、`gov`、`keys`、`oracle`、`random`、`service`、`slashing`、`staking`、`tendermint`。

## install

### Requirement

Go version above 1.13.5

### Use Go Mod

```text
require (
    github.com/irisnet/irishub-sdk-go latest
)
```

## Usage

### KeyDAO

在使用SDK之前，你需要实现管理Key的相关接口(`KeyDAO`)，主要是CRUD操作，接口定义如下：

```go
type KeyDAO interface {
    AccountAccess
    Crypto
}

type AccountAccess interface {
    Write(name string, store Store) error
    Read(name string) (Store, error)
    Delete(name string) error
}
type Crypto interface {
    Encrypt(data string, password string) (string, error)
    Decrypt(data string, password string) (string, error)
}
```

其中`Store`包括两种存储方式，一种是基于私钥的方式，定义如下：

```go
type KeyInfo struct {
    PrivKey string `json:"priv_key"`
    Address string `json:"address"`
}
```

另外一种是基于keystore的方式，定义如下：

```go
type KeystoreInfo struct {
    Keystore string `json:"keystore"`
}
```

你可以灵活选择其中任何一种私钥的管理方式。`Encrypt`和`Decrypt`接口是对key的加解密处理，如果用户不实现，将默认使用`AES`。示例如下：

`KeyDao`实现`AccountAccess`接口：

```go
type Memory map[string]types.Store

func (m Memory) Write(name string, store types.Store) error {
    m[name] = store
    return nil
}

func (m Memory) Read(name string) (types.Store, error) {
    return m[name], nil
}

func (m Memory) Delete(name string) error {
    delete(m, name)
    return nil
}
```

**注意**：如果你不使用发送交易的相关API，可以不实现`KeyDAO`接口。

### Init Client

实现`KeyDAO`接口之后，还需要配置`SDK`的一些参数，说明如下：

| 配置项    | 类型          | 描述                                               |
| --------- | ------------- | -------------------------------------------------- |
| NodeURI   | 字符串        | SDK连接的irishub节点RPC地址，例如：localhost:26657 |
| Network   | enum          | irishub网络类型，取值：`Testnet`、`Mainnet`        |
| ChainID   | string        | irishub的ChainID，例如：`irishub`                  |
| Gas       | uint64        | 交易需支付的最大Gas，例如：`20000`                 |
| Fee       | DecCoins      | 交易需支付的交易费                                 |
| KeyDAO    | KeyDAO        | 私钥管理接口                                       |
| Mode      | enum          | 交易的广播模式，取值：`Sync`、`Async`、`Commit`    |
| StoreType | enum          | 私钥的存储方式，取值：`Keystore`、`Key`、          |
| Timeout   | time.Duration | 交易的超时时间，例如：`5s`                         |
| Level     | string        | 日志输出级别，例如：`info`                         |

初始化`SDK`代码如下：


```go
client := sdk.NewClient(types.ClientConfig{
    NodeURI:   NodeURI,
    Network:   Network,
    ChainID:   ChainID,
    Gas:       Gas,
    Fee:       fees,
    KeyDAO:    types.NewDefaultKeyDAO(&Memory{}),
    Mode:      Mode,
    StoreType: types.Key,
    Timeout:   10 * time.Second,
    Level:     "info",
})
```

使用`NewDefaultKeyDAO`方法初始化一个`KeyDAO`实例，将默认使用`AES`的加密方式。

如果你想使用`SDK`发送一笔转账交易，示例如下：

```go
coins, err := types.ParseDecCoins("0.1iris")
to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
baseTx := types.BaseTx{
    From:     "username",
    Gas:      20000,
    Memo:     "test",
    Mode:     types.Commit,
    Password: "password",
}

result, err := client.Bank().Send(to, coins, baseTx)
```

有关更多API使用说明文档，请查看[文档](https://pkg.go.dev/mod/github.com/irisnet/irishub-sdk-go)。
