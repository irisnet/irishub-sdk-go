package gov

import "github.com/irisnet/irishub-sdk-go/types"

func NewGovClient(ac types.AbstractClient) Client {
	return govClient{
		AbstractClient: ac,
	}
}
