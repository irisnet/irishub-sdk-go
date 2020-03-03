package types

type Module interface {
	RegisterCodec(cdc Codec)
	Name() string
}
