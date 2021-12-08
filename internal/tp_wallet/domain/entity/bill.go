package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tp_wallet/internal/tp_wallet/adapter/dto"
	"tp_wallet/internal/tp_wallet/domain/vo"
)

type Bill struct {
	Id              primitive.ObjectID     `json:"id,omitempty" bson:"_id"`                                      // 订单id
	NumericalOrder  uint64                 `json:"numerical_order" bson:"numerical_order"`                       // 流水号
	Uid             uint64                 `json:"uid,omitempty" bson:"uid,omitempty"`                           // uid
	Cid             uint64                 `json:"cid,omitempty" bson:"cid"`                                     // 渠道id
	Amount          string                 `json:"amount,omitempty" bson:"amount,omitempty"`                     // 转账金额
	Gas             vo.Gas                 `json:"gas,omitempty" bson:"gas,omitempty"`                           // 手续费
	BillType        dto.TransferType       `json:"bill_type,omitempty" bson:"bill_type,omitempty"`               // 账单类型
	BillStatus      dto.BillStatus         `json:"bill_status,omitempty" bson:"bill_status,omitempty"`           // 账单状态
	FromUid         uint64                 `json:"from_uid,omitempty" bson:"from_uid,omitempty"`                 // From uid
	FromAddr        string                 `json:"from_addr,omitempty" bson:"from_addr,omitempty"`               // From地址
	ToUid           uint64                 `json:"to_uid,omitempty" bson:"to_uid,omitempty"`                     // To uid
	ToAddr          string                 `json:"to_addr,omitempty" bson:"to_addr,omitempty"`                   // To地址
	Hash            string                 `json:"hash,omitempty" bson:"hash,omitempty"`                         // 交易哈希
	Remark          string                 `json:"remark,omitempty" bson:"remark,omitempty"`                     // 备注
	Times           int64                  `json:"times,omitempty" bson:"times,omitempty"`                       // 重试次数（异步订单会有）
	CreateTime      int64                  `json:"create_time,omitempty" bson:"create_time,omitempty"`           // 创建时间
	UpdateTime      int64                  `json:"update_time,omitempty" bson:"update_time,omitempty"`           // 修改时间
	BalanceRecord   vo.BalanceRecord       `json:"balance_record,omitempty" bson:"balance_record,omitempty"`     // 金额记录
	ContractRecord  vo.ContractRecord      `json:"contract_record,omitempty" bson:"contract_record,omitempty"`   // nft记录
	IsBalanceTrade  bool                   `json:"is_balance_trade,omitempty" bson:"is_balance_trade,omitempty"` // 是否涉及修改余额
	Version         uint64                 `json:"version,omitempty" bson:"version"`                             // 版本，每次操作会+1，常用于乐观锁
	LastTimeVersion int64                  `bson:"-" json:"-"`                                                   // 用于更新
	operation       map[string]interface{} `json:"-" bson:"-"`
}

func (b *Bill) CheckBillCreate() bool {
	if b.Uid == 0 || b.Cid == 0 || b.BillType == 0 || b.BillStatus == 0 || len(b.FromAddr) == 0 || len(b.ToAddr) == 0 || b.NumericalOrder <= 0 {
		return false
	}
	return true
}
