package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tp_wallet/internal/wallet/dto"
)

type NftOwner struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                  // 订单id
	NftToken      string             `json:"nft_token,omitempty" bson:"nft_token"`               // nft唯一标识
	GameToken     string             `json:"game_token,omitempty" bson:"game_token"`             // 游戏唯一标识(所有修改，创建以游戏唯一标识为准)
	CreateHash    string             `json:"create_hash,omitempty" bson:"create_hash,omitempty"` // 创建hash(用于判断)
	ContractToken string             `json:"contract_token,omitempty" bson:"contract_token"`     // 合约标识
	OwnerAddress  string             `json:"owner_address,omitempty" bson:"owner_address"`       // 拥有者地址
	CreateTime    int64              `json:"create_time,omitempty" bson:"create_time"`           // 创建时间
	UpdateTime    int64              `json:"update_time,omitempty" bson:"update_time"`           // 修改时间
	NftData
	Uid uint64 `json:"uid,omitempty" bson:"-"` // 拥有者id
}

func (b NftOwner) ToPb() *dto.NftInfo {
	return &dto.NftInfo{
		Uid:           b.Uid,
		OwnerAddress:  b.OwnerAddress,
		GameId:        b.GameName,
		NftGameToken:  b.GameToken,
		NftChainToken: b.NftToken,
		Level:         b.Level,
		ContractToken: b.ContractToken,
	}
}

type NftContract struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id"`                                    // 订单id
	NftToken           string             `json:"nft_token,omitempty" bson:"nft_token"`                       // nft唯一标识
	GameToken          string             `json:"game_token,omitempty" bson:"game_token"`                     // nft游戏内唯一标识
	ContractToken      string             `json:"contract_token,omitempty" bson:"contract_token"`             // 合约信息
	OwnerAddress       string             `json:"owner_address,omitempty" bson:"owner_address"`               // 拥有者地址
	ReallyOwnerAddress string             `json:"really_owner_address,omitempty" bson:"really_owner_address"` // 真实拥有者地址
	CreateTime         int64              `json:"create_time,omitempty" bson:"create_time"`                   // 创建时间
	UpdateTime         int64              `json:"update_time,omitempty" bson:"update_time"`                   // 修改时间
	DeleteTime         int64              `json:"delete_time,omitempty" bson:"delete_time"`                   // 删除时间
	NftData
	Uid uint64 `json:"uid,omitempty" bson:"-"` // 拥有者id
}

func (b NftContract) ToPb() *dto.NftInfo {
	return &dto.NftInfo{
		Uid:           b.Uid,
		OwnerAddress:  b.OwnerAddress,
		GameId:        b.GameName,
		NftGameToken:  b.GameToken,
		NftChainToken: b.NftToken,
		Level:         b.Level,
		ContractToken: b.GameName,
	}
}

type NftChangeRecord struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id"`                        // 订单id
	NftToken      string             `json:"nft_token,omitempty" bson:"nft_token"`           // nft唯一标识
	GameToken     string             `json:"game_token,omitempty" bson:"nft_token"`          // nft唯一标识
	Hash          string             `json:"hash,omitempty" bson:"hash"`                     // 交易哈希
	BeforeAddress string             `json:"before_address,omitempty" bson:"before_address"` // 变更之前地址
	AfterAddress  string             `json:"after_address,omitempty" bson:"after_address"`   // 变更之后地址
	CreateTime    int64              `json:"create_time,omitempty" bson:"create_time"`       // 创建时间
}
