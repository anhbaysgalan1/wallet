package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddressPrivate struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`                 // id
	Address    string             `json:"address" bson:"address"`         // 地址
	Private    string             `json:"private" bson:"private"`         // 私钥
	Currency   string             `json:"currency" bson:"currency"`       // 币种
	Net        string             `json:"net" bson:"net"`                 // 网络
	Status     AddressStatus      `json:"status" bson:"status"`           // 状态
	CreateTime int64              `json:"create_time" bson:"create_time"` // 创建时间
	UpdateTime int64              `json:"update_time" bson:"update_time"` // 更新时间
	Remarks    string             `json:"remarks" bson:"remarks"`         // 备注
}

type AddressStatus uint16

const (
	AddressStatus_UnUsed AddressStatus = 1 // 未使用
	AddressStatus_Used   AddressStatus = 2 // 已使用
	AddressStatus_Admin  AddressStatus = 3 // 管理员
)
