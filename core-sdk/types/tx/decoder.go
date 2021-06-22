package tx

import (
	"github.com/irisnet/irishub-sdk-go/common/codec"
	"github.com/irisnet/irishub-sdk-go/common/codec/unknownproto"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// DefaultTxDecoder returns a default protobuf TxDecoder using the provided Marshaler.
func DefaultTxDecoder(cdc *codec.ProtoCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var raw TxRaw

		// reject all unknown proto fields in the root TxRaw
		err := unknownproto.RejectUnknownFieldsStrict(txBytes, &raw)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(txBytes, &raw)
		if err != nil {
			return nil, err
		}

		var body TxBody

		// allow non-critical unknown fields in TxBody
		txBodyHasUnknownNonCriticals, err := unknownproto.RejectUnknownFields(raw.BodyBytes, &body, true)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(raw.BodyBytes, &body)
		if err != nil {
			return nil, err
		}

		var authInfo AuthInfo

		// reject all unknown proto fields in AuthInfo
		err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, err
		}

		err = cdc.UnmarshalBinaryBare(raw.AuthInfoBytes, &authInfo)
		if err != nil {
			return nil, err
		}

		theTx := &Tx{
			Body:       &body,
			AuthInfo:   &authInfo,
			Signatures: raw.Signatures,
		}

		return &wrapper{
			tx:                           theTx,
			bodyBz:                       raw.BodyBytes,
			authInfoBz:                   raw.AuthInfoBytes,
			txBodyHasUnknownNonCriticals: txBodyHasUnknownNonCriticals,
		}, nil
	}
}

// DefaultJSONTxDecoder returns a default protobuf JSON TxDecoder using the provided Marshaler.
func DefaultJSONTxDecoder(cdc *codec.ProtoCodec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var theTx Tx
		err := cdc.UnmarshalJSON(txBytes, &theTx)
		if err != nil {
			return nil, err
		}

		return &wrapper{
			tx: &theTx,
		}, nil
	}
}
