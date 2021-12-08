package vo

import (
	"encoding/json"
	"tp_wallet/internal/tp_wallet/adapter/dto"
)

type ErrorBillRemark struct {
	From              uint64           `json:"from,omitempty" bson:"from,omitempty"`                               // 转出uid
	FromBeforeBalance string           `json:"from_before_balance,omitempty" bson:"from_before_balance,omitempty"` // 转出之前金额
	FromAfterBalance  string           `json:"from_after_balance,omitempty" bson:"from_after_balance,omitempty"`   // 转出之后金额
	To                uint64           `json:"to,omitempty" bson:"from,omitempty"`                                 // 转出uid
	ToBeforeBalance   string           `json:"to_before_balance,omitempty" bson:"to_before_balance,omitempty"`     // 转出之前金额
	ToAfterBalance    string           `json:"to_after_balance,omitempty" bson:"to_after_balance,omitempty"`       // 转出之后金额
	ContractType      dto.ContractType `json:"contract_type,omitempty" bson:"contract_type"`                       // 合约类型
	NftToken          string           `json:"nft_token,omitempty" bson:"nft_token"`                               // nft token
	GameToken         string           `json:"game_token,omitempty" bson:"game_token"`                             // game token
	Data              string           `json:"data,omitempty" bson:"data,omitempty"`                               // 错误信息
}

func (e ErrorBillRemark) ToJson() string {
	result, _ := json.Marshal(e)
	return string(result)
}

func (e *ErrorBillRemark) UnJson(req string) error {
	return json.Unmarshal([]byte(req), e)
}
