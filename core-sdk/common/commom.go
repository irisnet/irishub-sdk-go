package common

import (
	"crypto/rand"
	"github.com/irisnet/irishub-sdk-go/types"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func SubArray(subLen int, array types.SplitAble) (segments []types.SplitAble) {
	maxLen := array.Len()
	if maxLen <= subLen {
		return []types.SplitAble{array}
	}

	batch := maxLen / subLen
	if maxLen%subLen > 0 {
		batch++
	}

	for i := 1; i <= batch; i++ {
		start := (i - 1) * subLen
		end := i * subLen
		if i != batch {
			segments = append(segments, array.Sub(start, end))
		} else {
			segments = append(segments, array.Sub(start, array.Len()))
		}
	}
	return segments
}

func ParsePage(page uint64, size uint64) (offset, limit uint64) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	offset = (page - 1) * size
	limit = size
	return
}
