package tools

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitArray(t *testing.T) {
	data := Ints{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	subData := SplitArray(4, data)
	require.Len(t, subData, 3)
}

type Ints []int

func (i Ints) Len() int {
	return len(i)
}

func (i Ints) Sub(begin, end int) sdk.SplitAble {
	return i[begin:end]
}
