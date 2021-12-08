package vo

type Gas struct {
	GasPrice    string `json:"gas_price,omitempty" bson:"gas_price"`
	GasLimit    string `json:"gas_limit,omitempty" bson:"gas_limit"`
	GasCurrency string `json:"gas_currency,omitempty" bson:"gas_currency,omitempty"` // 手续费币种
}

func (nd Gas) IsEmpty() bool {
	return nd == Gas{}
}
