package types

type Module interface {
	RegisterCodec(cdc Codec)
	//RegisterErrorCode()
	Name() string
}

//The purpose of this interface is to convert the irishub system type to the user receiving type
// and standardize the user interface
type Response interface {
	Convert() interface{}
}
