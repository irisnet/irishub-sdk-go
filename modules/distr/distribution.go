package distr

import "github.com/irisnet/irishub-sdk-go/types"

func NewClient(ac types.AbstractClient) Client {
	return distrClient{
		AbstractClient: ac,
	}
}
