package staking

import (
	yaml "gopkg.in/yaml.v2"
)

// String implements the stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// String implements the Stringer interface for a Commission object.
func (c Commission) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

// String implements the Stringer interface for a CommissionRates object.
func (cr CommissionRates) String() string {
	out, _ := yaml.Marshal(cr)
	return string(out)
}

// String implements the Stringer interface for a Description object.
func (d Description) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// String implements the Stringer interface for a Validator object.
func (v Validator) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}
