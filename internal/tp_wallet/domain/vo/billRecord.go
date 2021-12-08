package vo

type BalanceRecord struct {
	BeforeBalance string `json:"before_balance,omitempty" bson:"before_balance"` // 转出之前金额
	AfterBalance  string `json:"after_balance,omitempty" bson:"after_balance"`   // 转出之后金额
}

func (nd BalanceRecord) IsEmpty() bool {
	return nd == BalanceRecord{}
}

type ContractRecord struct {
	ContractType string `json:"contract_type,omitempty" bson:"contract_type"` // 合约类型
	ContractAddr string `json:"contract_addr,omitempty" bson:"contract_addr"` // 合约地址
	NftToken     string `json:"nft_token,omitempty" bson:"nft_token"`         // nft 类型
	GameToken    string `json:"game_token,omitempty" bson:"game_token"`       // nft 类型
}

func (nd ContractRecord) IsEmpty() bool {
	return nd == ContractRecord{}
}
