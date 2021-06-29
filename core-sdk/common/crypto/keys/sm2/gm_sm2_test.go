package sm2

import (
	"fmt"
	"testing"
)

func TestSm2(t *testing.T) {
	Sm2()
}

func TestGenerateKey(t *testing.T) {
	priv := GenerateKey()

	fmt.Println(priv.Public())
	fmt.Println(GetPublickey(priv))
}
