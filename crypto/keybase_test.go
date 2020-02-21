package crypto

import (
	"fmt"
	"testing"
)

func TestKeyManager_ExportAsPrivateKey(t *testing.T) {
	mnemonic := "refuse salmon will muscle exclude artist cancel hunt such brand latin collect tongue train saddle scorpion transfer mass loop earth settle between camp slush"
	keyManager, _ := NewMnemonicKeyManager(mnemonic)
	priKey, _ := keyManager.ExportAsPrivateKey()
	fmt.Println(priKey)
}
