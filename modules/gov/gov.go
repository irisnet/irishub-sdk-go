package gov

import "github.com/irisnet/irishub-sdk-go/types"

func NewGovClient(ac types.AbstractClient) Gov {
	return govClient{
		AbstractClient: ac,
	}
}
