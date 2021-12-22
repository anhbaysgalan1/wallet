package entity

import (
	"encoding/json"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bill struct {
	Id             primitive.ObjectID    `json:"id,omitempty" bson:"_id"`                          // 订单id
	NumericalOrder uint64                `json:"numerical_order" bson:"numerical_order"`           // 流水号
	Uid            uint64                `json:"uid,omitempty" bson:"uid"`                         // uid
	Cid            uint64                `json:"cid,omitempty" bson:"cid"`                         // 渠道id
	Gas            Gas                   `json:"gas,omitempty" bson:"gas"`                         // 手续费
	BillType       BillType              `json:"bill_type,omitempty" bson:"bill_type"`             // 账单类型
	TransferType   walletPb.TransferType `json:"transfer_type,omitempty" bson:"transfer_type"`     // 交易类型
	BillStatus     walletPb.BillStatus   `json:"bill_status,omitempty" bson:"bill_status"`         // 账单状态
	FromUid        uint64                `json:"from_uid,omitempty" bson:"from_uid"`               // From uid
	FromAddr       string                `json:"from_addr,omitempty" bson:"from_addr"`             // From地址
	ToUid          uint64                `json:"to_uid,omitempty" bson:"to_uid"`                   // To uid
	ToAddr         string                `json:"to_addr,omitempty" bson:"to_addr"`                 // To地址
	Hash           string                `json:"hash,omitempty" bson:"hash"`                       // 交易哈希
	Remark         string                `json:"remark,omitempty" bson:"remark"`                   // 备注
	Times          int64                 `json:"times,omitempty" bson:"times"`                     // 重试次数（异步订单会有）
	CreateTime     int64                 `json:"create_time,omitempty" bson:"create_time"`         // 创建时间
	UpdateTime     int64                 `json:"update_time,omitempty" bson:"update_time"`         // 修改时间
	BalanceRecord  BalanceRecord         `json:"balance_record,omitempty" bson:"balance_record"`   // 金额记录
	ContractRecord ContractRecord        `json:"contract_record,omitempty" bson:"contract_record"` // nft记录
	TpRecord       TpRecord              `json:"tp_record,omitempty" bson:"tp_record"`             // 第三方透传数据
}

func (b *Bill) CheckBillCreate() bool {
	if b.Uid == 0 || b.TransferType == 0 || b.BillType == 0 || b.BillStatus == 0 || b.NumericalOrder <= 0 || b.FromUid == 0 || b.ToUid == 0 {
		return false
	}
	return true
}

func (b *Bill) DealWithBillToCash(billForStore *Bill) {
	b.Hash = ""
	b.FromUid = 0
	b.FromAddr = ""
	b.ToUid = 0
	b.ToAddr = ""
	b.BalanceRecord.Amount = billForStore.BalanceRecord.Amount
	b.BalanceRecord.BeforeBalance = billForStore.BalanceRecord.BeforeBalance
	b.BalanceRecord.AfterBalance = billForStore.BalanceRecord.AfterBalance
	b.BalanceRecord.Currency = billForStore.BalanceRecord.Currency
}

func (b *Bill) DealWithBillPending() bool {
	b.Id = primitive.NilObjectID
	b.Uid = 0
	b.Cid = 0
	b.BillType = 0
	b.TransferType = 0
	b.FromUid = 0
	b.ToUid = 0
	return true
}

func (b Bill) ToPb() *walletPb.BillInfo {
	var result = &walletPb.BillInfo{
		Id:             b.Id.Hex(),
		NumericalOrder: b.NumericalOrder,
		Uid:            b.Uid,
		TransferType:   b.TransferType,
		BillStatus:     b.BillStatus,
		Hash:           b.Hash,
		FromUid:        b.FromUid,
		FromAddr:       b.FromAddr,
		ToUid:          b.ToUid,
		ToAddr:         b.ToAddr,
		Remark:         b.Remark,
		Times:          b.Times,
		BillType:       int64(b.BillType),
	}
	if !b.Gas.IsEmpty() {
		result.Gas = b.Gas.Gas
		result.GasCurrency = b.Gas.GasCurrency
	}
	if !b.BalanceRecord.IsEmpty() {
		result.BalanceRecord = &walletPb.BalanceRecord{
			Currency:      b.BalanceRecord.Currency,
			Amount:        b.BalanceRecord.Amount,
			ReceiveAmount: b.BalanceRecord.ReceiveAmount,
			BeforeBalance: string(b.BalanceRecord.BeforeBalance),
			AfterBalance:  string(b.BalanceRecord.AfterBalance),
		}
	}
	if !b.ContractRecord.IsEmpty() {
		result.ContractRecord = &walletPb.ContractRecord{
			ContractType: b.ContractRecord.ContractType,
			ContractAddr: b.ContractRecord.ContractAddr,
			NftToken:     b.ContractRecord.NftToken,
			GameToken:    b.ContractRecord.GameToken,
			Num:          b.ContractRecord.Num,
		}
	}
	if !b.TpRecord.IsEmpty() {
		result.TpRecord = &walletPb.TpRecord{
			OrderId: b.TpRecord.OrderId,
			Type:    b.TpRecord.Type,
			Remarks: b.TpRecord.Remarks,
			Data:    b.TpRecord.Data,
		}
	}
	return result
}

type BillType uint16

const (
	BillType_Eip20   BillType = 1 // eip20代币
	BillType_Eip721  BillType = 2 // eip721代币
	BillType_Eip1155 BillType = 3 // eip1155代币
)

type Gas struct {
	Gas         string `json:"gas,omitempty" bson:"gas"` // 手续费
	GasPrice    string `json:"gas_price,omitempty" bson:"gas_price"`
	GasLimit    string `json:"gas_limit,omitempty" bson:"gas_limit"`
	GasCurrency string `json:"gas_currency,omitempty" bson:"gas_currency"` // 手续费币种
}

func (nd Gas) IsEmpty() bool {
	return nd == Gas{}
}

type BalanceRecord struct {
	Currency      string `json:"currency,omitempty" bson:"currency"` // 交易币种
	Amount        string `json:"amount,omitempty" bson:"amount"`     // 转账金额
	ReceiveAmount string `json:"receive_amount,omitempty" bson:"receive_amount"`
	BeforeBalance Amount `json:"before_balance,omitempty" bson:"before_balance"` // 转出之前金额
	AfterBalance  Amount `json:"after_balance,omitempty" bson:"after_balance"`   // 转出之后金额
}

func (nd BalanceRecord) IsEmpty() bool {
	return nd == BalanceRecord{}
}

type ContractRecord struct {
	ContractType string `json:"contract_type,omitempty" bson:"contract_type"` // 合约类型
	ContractAddr string `json:"contract_addr,omitempty" bson:"contract_addr"` // 合约地址
	GameId       string `json:"game_id,omitempty" bson:"game_id"`             // 游戏唯一标识
	NftToken     string `json:"nft_token,omitempty" bson:"nft_token"`         // 道具在合约内唯一标识
	GameToken    string `json:"game_token,omitempty" bson:"game_token"`       // 道具在游戏内唯一标识
	Num          uint64 `json:"num,omitempty" bson:"num"`                     // 数量
}

func (nd ContractRecord) IsEmpty() bool {
	return nd == ContractRecord{}
}

type TpRecord struct {
	OrderId string `json:"order_id,omitempty" bson:"order_id"` // 第三方订单id
	Type    uint32 `json:"type,omitempty" bson:"type"`         // 第三方转账类型
	Remarks string `json:"remarks,omitempty" bson:"remarks"`   // 备注
	Data    string `json:"data,omitempty" bson:"data"`         // 透传数据
}

func (nd TpRecord) IsEmpty() bool {
	return nd == TpRecord{}
}

type ErrorBillRemark struct {
	From              uint64 `json:"from,omitempty" bson:"from"`                               // 转出uid
	FromBeforeBalance string `json:"from_before_balance,omitempty" bson:"from_before_balance"` // 转出之前金额
	FromAfterBalance  string `json:"from_after_balance,omitempty" bson:"from_after_balance"`   // 转出之后金额
	To                uint64 `json:"to,omitempty" bson:"from"`                                 // 转出uid
	ToBeforeBalance   string `json:"to_before_balance,omitempty" bson:"to_before_balance"`     // 转出之前金额
	ToAfterBalance    string `json:"to_after_balance,omitempty" bson:"to_after_balance"`       // 转出之后金额
	Data              string `json:"data,omitempty" bson:"data"`                               // 错误信息
}

func (e ErrorBillRemark) ToJson() string {
	result, _ := json.Marshal(e)
	return string(result)
}

func (e *ErrorBillRemark) UnJson(req string) error {
	return json.Unmarshal([]byte(req), e)
}

func PbToBill(info *walletPb.BillInfo) *Bill {
	if info == nil {
		return nil
	}
	bill := &Bill{
		NumericalOrder: info.NumericalOrder,
		Uid:            info.Uid,
		BillType:       BillType(info.GetBillType()),
		TransferType:   info.TransferType,
		BillStatus:     info.BillStatus,
		FromUid:        info.GetFromUid(),
		FromAddr:       info.GetFromAddr(),
		ToUid:          info.ToUid,
		ToAddr:         info.GetToAddr(),
		Hash:           info.GetHash(),
		Remark:         info.GetRemark(),
		Times:          info.GetTimes(),
	}
	id, err := primitive.ObjectIDFromHex(info.Id)
	if err == nil {
		bill.Id = id
	}
	if info.Gas != "" && info.GasCurrency != "" {
		bill.Gas = Gas{
			Gas:         info.Gas,
			GasPrice:    "",
			GasLimit:    "",
			GasCurrency: info.GasCurrency,
		}
	}
	if info.GetBalanceRecord() != nil {
		bill.BalanceRecord = BalanceRecord{
			Currency:      info.BalanceRecord.Currency,
			Amount:        info.BalanceRecord.Amount,
			ReceiveAmount: info.BalanceRecord.ReceiveAmount,
		}
	}
	if info.GetContractRecord() != nil {
		bill.ContractRecord = ContractRecord{
			ContractType: info.ContractRecord.ContractType,
			ContractAddr: info.ContractRecord.ContractAddr,
			NftToken:     info.ContractRecord.NftToken,
			GameToken:    info.ContractRecord.GameToken,
			Num:          info.ContractRecord.Num,
		}
	}
	if info.GetTpRecord() != nil {
		bill.TpRecord = TpRecord{
			OrderId: info.TpRecord.OrderId,
			Type:    info.TpRecord.Type,
			Remarks: info.TpRecord.Remarks,
			Data:    info.TpRecord.Data,
		}
	}
	return bill
}
