package json

import "encoding/json"

// MustSort is like Sort but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSort(toSortJSON []byte) []byte {
	js, err := Sort(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}

// SortedJSON takes any JSON and returns it sorted by keys. Also, all white-spaces
// are removed.
// This method can be used to canonicalize JSON to be returned by GetSignBytes,
// e.g. for the ledger integration.
// If the passed JSON isn't valid it will return an error.
func Sort(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}
