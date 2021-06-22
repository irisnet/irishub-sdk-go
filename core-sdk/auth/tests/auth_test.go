package tests

//functional test
import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/auth"
	"testing"
)

const (
	nodeURI  = "tcp://localhost:26657"
	grpcAddr = "localhost:9090"
	chainID  = "test"
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	addr     = "iaa1w9lvhwlvkwqvg08q84n2k4nn896u9pqx93velx"
)

func Test(t *testing.T) {

}
func TestBaseAccountGetAddress(t *testing.T) {
	acc := auth.BaseAccount{}
	fmt.Println(acc.GetAddress())
}

func TestGo(t *testing.T) {

}
