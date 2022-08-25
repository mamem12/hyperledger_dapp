package model

type Approval struct {
	Spender   string `json:"spender"`
	Owner     string `json:"Owner"`
	Allowance int    `json:"allowance"`
}

func NewApproval(spender, owner string, allowance int) *Approval {

	return &Approval{
		Spender:   spender,
		Owner:     owner,
		Allowance: allowance,
	}

}
