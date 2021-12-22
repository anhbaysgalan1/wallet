package entity

import (
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/big"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

type Balance struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Uid        uint64             `json:"uid,omitempty" bson:"uid"`
	Cid        uint64             `json:"cid,omitempty" bson:"cid"`           // 渠道id
	Currency   string             `json:"currency,omitempty" bson:"currency"` // 币种
	Balance    Amount             `json:"balance,omitempty" bson:"balance"`   // 余额
	CreateTime int64              `json:"create_time,omitempty" bson:"create_time"`
	UpdateTime int64              `json:"update_time,omitempty" bson:"update_time"`
}

func (b *Balance) CheckCreate() bool {
	if b.Uid == 0 || len(b.Currency) == 0 {
		return false
	}
	now := tool.GetTimeUnixMilli()
	b.CreateTime = now
	b.UpdateTime = now
	return true
}

type Amount string

var bigZero = new(big.Int)

func (a Amount) SetByNum(num uint64) Amount {
	if a == "" {
		var str = Amount("0")
		a = str
	}
	var str = new(big.Int).SetUint64(num).String()
	var amount = Amount(str)
	return amount
}

func (a Amount) SetByStr(num string) (Amount, error) {
	if a == "" {
		var str = Amount("0")
		a = str
	}
	var result bool
	if _, result = new(big.Int).SetString(num, 0); !result {
		return "", hcode.ErrInternalParameter
	}
	var amount = Amount(num)
	return amount, nil
}

func (a Amount) Cmp(num string) (int, error) {
	if a == "" {
		var str = Amount("0")
		a = str
	}
	var amount, result = new(big.Int).SetString(string(a), 0)
	if !result {
		return 0, hcode.ErrInternalParameter
	}
	var cmp = new(big.Int)
	cmp, result = cmp.SetString(num, 0)
	if !result {
		return 0, hcode.ErrInternalParameter
	}
	return amount.Cmp(cmp), nil
}

func (a Amount) Add(num string) (Amount, error) {
	if a == "" {
		var str = Amount("0")
		a = str
	}
	var result bool
	var before, after, value = new(big.Int), new(big.Int), new(big.Int)
	_, result = value.SetString(num, 0)
	if !result {
		log.GetLogger().Error("[amount.Add] failed", zap.Any("amount", a), zap.String("num", num))
		return "", hcode.ErrBalanceConversion
	}
	if value.Cmp(bigZero) == 0 { // 加0，直接返回
		return "", nil
	}
	if value.Cmp(bigZero) < 0 { // 不能加负数
		log.GetLogger().Error("[amount.Add] failed", zap.Any("amount", a), zap.String("num", num))
		return "", hcode.ErrBalanceMinus
	}
	_, result = before.SetString(string(a), 0)
	if !result {
		log.GetLogger().Error("[amount.Add] failed", zap.Any("amount", a), zap.String("num", num))
		return "", hcode.ErrBalanceConversion
	}
	after = before.Add(before, value)
	return Amount(after.String()), nil
}

func (a *Amount) Sub(num string) (Amount, bool, error) { // 返回结果是否为正数
	if *a == "" {
		var str = Amount("0")
		a = &str
	}
	var result bool
	var before, after, value = new(big.Int), new(big.Int), new(big.Int)
	_, result = value.SetString(num, 0)
	if !result {
		log.GetLogger().Error("[amount.Sub] failed", zap.Any("amount", a), zap.String("num", num))
		return "", false, hcode.ErrBalanceConversion
	}
	if value.Cmp(bigZero) == 0 { // 加0，直接返回
		return "", true, nil
	}
	if value.Cmp(bigZero) < 0 { // 不能减负数
		log.GetLogger().Error("[amount.Sub] failed", zap.Any("amount", a), zap.String("num", num))
		return "", false, hcode.ErrBalanceMinus
	}
	_, result = before.SetString(string(*a), 0)
	if !result {
		log.GetLogger().Error("[amount.Sub] failed", zap.Any("amount", a), zap.String("num", num))
		return "", false, hcode.ErrBalanceConversion
	}
	after = before.Sub(before, value)
	if after.Cmp(bigZero) < 0 {
		return "", false, nil
	} else {
		return Amount(after.String()), true, nil
	}
}

func AmountMapToString(a map[string]Amount) map[string]string {
	var result = make(map[string]string)
	for k, v := range a {
		result[k] = string(v)
	}
	return result
}
