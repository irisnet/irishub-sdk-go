package gov

import yaml "gopkg.in/yaml.v2"

func (d Deposit) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

func (tr TallyResult) String() string {
	out, _ := yaml.Marshal(tr)
	return string(out)
}

func (v Vote) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

func (dp DepositParams) String() string {
	out, _ := yaml.Marshal(dp)
	return string(out)
}

func (vp VotingParams) String() string {
	out, _ := yaml.Marshal(vp)
	return string(out)
}

func (tp TallyParams) String() string {
	out, _ := yaml.Marshal(tp)
	return string(out)
}
