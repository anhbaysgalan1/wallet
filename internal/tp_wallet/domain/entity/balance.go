package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/tool"
)

var bigZero = new(big.Int)

type Balance struct {
	Id              primitive.ObjectID     `json:"id,omitempty" bson:"_id"`
	Cid             uint64                 `json:"cid,omitempty" bson:"cid"`           // 渠道id
	Uid             uint64                 `json:"uid,omitempty" bson:"uid"`           // 用户id
	Currency        string                 `json:"currency,omitempty" bson:"currency"` // 币种
	Balance         Amount                 `json:"balance,omitempty" bson:"balance"`   // 余额
	CreateTime      int64                  `json:"create_time,omitempty" bson:"create_time"`
	UpdateTime      int64                  `json:"update_time,omitempty" bson:"update_time"`
	Version         uint64                 `json:"version,omitempty" bson:"version"` // 版本，每次操作会+1，常用于乐观锁
	LastTimeVersion int64                  `bson:"-" json:"-"`                       // 用于更新
	operation       map[string]interface{} `json:"-" bson:"-"`
}

func (b *Balance) CheckCreate() bool {
	if b.Uid == 0 || len(b.Currency) == 0 || b.Cid == 0 {
		return false
	}
	now := tool.GetTimeUnixMilli()
	b.Balance.SetByNum(0)
	b.CreateTime = now
	b.UpdateTime = now
	b.Version = 1
	return true
}

func (b *Balance) GetUpdates() (filter map[string]interface{}, updates map[string]interface{}) {
	if b.operation == nil {
		return nil, nil
	}
	updates = make(map[string]interface{})
	filter = make(map[string]interface{})
	filter["_id"] = b.Id
	for k, v := range b.operation {
		updates[k] = v
	}
	if b.LastTimeVersion != 0 {
		filter["version"] = b.LastTimeVersion
		b.LastTimeVersion = 0
	}
	b.operation = nil
	return
}

func (b *Balance) setOpt(name string, value interface{}) {
	if b.operation == nil {
		b.operation = make(map[string]interface{})
	}
	b.operation[name] = value
}

func (b *Balance) SetBalanceSetByNum(num uint64) {
	b.Balance.SetByNum(num)
	b.setOpt("balance", b.Balance)
}

func (b *Balance) SetBalanceByStr(num string) {
	b.Balance.SetByStr(num)
	b.setOpt("balance", b.Balance)
}

func (b *Balance) SetBalanceAdd(num string) (err error) {
	err = b.Balance.Add(num)
	if err != nil {
		return err
	}
	b.setOpt("balance", b.Balance)
	return nil
}

func (b *Balance) SetBalanceSub(num string) (result bool, err error) {
	result, err = b.Balance.Sub(num)
	if err != nil {
		return result, err
	}
	b.setOpt("balance", b.Balance)
	return result, err
}

type Amount string

func (a *Amount) SetByNum(num uint64) {
	if *a == "" {
		var str = Amount("0")
		a = &str
	}
	var str = new(big.Int).SetUint64(num).String()
	var amount = Amount(str)
	a = &amount
}

func (a *Amount) SetByStr(num string) {
	if *a == "" {
		var str = Amount("0")
		a = &str
	}
	var amount = Amount(num)
	a = &amount
}

func (a *Amount) Add(num string) error {
	if *a == "" {
		var str = Amount("0")
		a = &str
	}
	var result bool
	var before, after, value = new(big.Int), new(big.Int), new(big.Int)
	_, result = value.SetString(num, 10)
	if !result {
		return hcode.ErrBalanceConversion
	}
	if value.Cmp(bigZero) == 0 { // 加0，直接返回
		return nil
	}
	if value.Cmp(bigZero) < 0 { // 不能加负数
		return hcode.ErrBalanceMinus
	}
	_, result = before.SetString(string(*a), 10)
	if !result {
		return hcode.ErrBalanceConversion
	}
	after = before.Add(before, value)
	a.SetByStr(after.String())
	return nil
}

func (a *Amount) Sub(num string) (bool, error) { // 返回结果是否为正数
	if *a == "" {
		var str = Amount("0")
		a = &str
	}
	var result bool
	var before, after, value = new(big.Int), new(big.Int), new(big.Int)
	_, result = value.SetString(num, 10)
	if !result {
		return false, hcode.ErrBalanceConversion
	}
	if value.Cmp(bigZero) == 0 { // 加0，直接返回
		return true, nil
	}
	if value.Cmp(bigZero) < 0 { // 不能减负数
		return false, hcode.ErrBalanceMinus
	}
	_, result = before.SetString(string(*a), 10)
	if !result {
		return false, hcode.ErrBalanceConversion
	}
	after = before.Sub(before, value)
	a.SetByStr(after.String())
	if after.Cmp(bigZero) < 0 {
		return false, nil
	} else {
		return true, nil
	}
}
