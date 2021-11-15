package entity

import (
	"github.com/leaf-rain/wallet/internal/account/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EntityAddressPrivate struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`                 // id
	Uid        string             `json:"uid" bson:"uid"`                 // uid
	Address    string             `json:"address" bson:"address"`         // 地址
	Private    string             `json:"private" bson:"private"`         // 私钥
	Currency   string             `json:"currency" bson:"currency"`       // 币种
	Status     dto.AccountType    `json:"status" bson:"status"`           // 状态
	CreateTime int64              `json:"create_time" bson:"create_time"` // 创建时间
	UpdateTime int64              `json:"update_time" bson:"update_time"` // 更新时间
	Remarks    string             `json:"remarks" bson:"remarks"`         // 备注
}
