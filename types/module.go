package types

type Module interface {
	RegisterCodec(cdc Codec)
	Name() string
}

//The purpose of this interface is to convert the irishub system type to the user receiving type
// and standardize the user interface
type Response interface {
	Convert() interface{}
}

type SplitAble interface {
	Len() int
	Sub(begin, end int) SplitAble
}
