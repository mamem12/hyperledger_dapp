package model

type Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"owner"`
	TotalSupply uint64 `json:"totalSupply"`
}

func NewMetadata(name, symbol, owner string, totalSupply uint) *Metadata {

	return &Metadata{
		Name:        name,
		Symbol:      symbol,
		Owner:       owner,
		TotalSupply: uint64(totalSupply),
	}
}
