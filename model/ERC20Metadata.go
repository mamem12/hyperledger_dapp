package model

type ERC20Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"owner"`
	TotalSupply uint64 `json:"totalSupply"`
}

func NewERC20Metadata(name, symbol, owner string, totalSupply uint) *ERC20Metadata {

	return &ERC20Metadata{
		Name:        name,
		Symbol:      symbol,
		Owner:       owner,
		TotalSupply: uint64(totalSupply),
	}
}

func (erc20 *ERC20Metadata) GetTotalSupply() *uint64 {
	return &erc20.TotalSupply
}

func (erc20 *ERC20Metadata) GetName() *string {
	return &erc20.Name
}

func (erc20 *ERC20Metadata) GetSymbol() *string {
	return &erc20.Symbol
}

func (erc20 *ERC20Metadata) GetOwner() *string {
	return &erc20.Owner
}
