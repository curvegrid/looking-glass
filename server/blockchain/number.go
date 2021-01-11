package blockchain

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Number represents an input number that might be inputted as either int or string.
type Number struct {
	*big.Int
}

var (
	_ json.Unmarshaler = (*Number)(nil)
)

// UnmarshalInt tries to unmarshal an integer value either from
// a JSON string or number.
func unmarshalInt(b []byte) (n *big.Int, err error) {
	if err = json.Unmarshal(b, &n); err == nil {
		return n, nil
	}
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	n, ok := n.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("cannot parse \"%s\" as int", s)
	}
	return
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Number) UnmarshalJSON(b []byte) error {
	v, err := unmarshalInt(b)
	n.Int = v
	return err
}

// GetInt returns the int. Covers the null case.
func (n *Number) GetInt() *big.Int {
	if n == nil {
		return nil
	}
	return n.Int
}
