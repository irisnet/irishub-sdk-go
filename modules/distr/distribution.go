package distr

import "github.com/irisnet/irishub-sdk-go/types"

func NewDistrClient(ac types.AbstractClient) Client {
	return distrClient{
		AbstractClient: ac,
	}
}
