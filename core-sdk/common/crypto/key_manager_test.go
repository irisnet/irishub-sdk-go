package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/irisnet/irishub-sdk-go/common/crypto"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

func TestNewMnemonicKeyManager(t *testing.T) {
	mnemonic := "nerve leader thank marriage spice task van start piece crowd run hospital control outside cousin romance left choice poet wagon rude climb leisure spring"

	km, err := crypto.NewMnemonicKeyManager(mnemonic, "secp256k1")
	assert.NoError(t, err)

	pubKey := km.ExportPubKey()
	pubkeyBech32, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey)
	assert.NoError(t, err)
	assert.Equal(t, "iap1qf6rwt2vpsdx9tcwq03w4dw9udd657u0gmknjd4l0ht699x6npll6hf0ru9", pubkeyBech32)

	address := sdk.AccAddress(pubKey.Address()).String()
	assert.Equal(t, "iaa1y9kd9uy7a4qnjp0z5yjx5jhrkv2ycdkzqc0h8z", address)
}
