package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id"` // id
	Uid         uint64             `json:"uid,omitempty" bson:"uid"`
	Cid         uint64             `json:"cid,omitempty" bson:"cid"`                   // 渠道id
	Currency    string             `json:"currency,omitempty" bson:"currency"`         // 币种
	Address     string             `json:"address,omitempty" bson:"address"`           // 地址
	AccountType AccountType        `json:"account_type,omitempty" bson:"account_type"` // 1 系统用户，2 普通用户
}

func (a Address) CheckAddressCreate() bool {
	if a.Uid <= 0 || len(a.Address) <= 0 || len(a.Currency) <= 0 {
		return false
	}
	return true
}

type AccountType uint16

const (
	AccountType_Admin    = 1
	AccountType_Ordinary = 2
	AccountType_Reboot   = 3
)
