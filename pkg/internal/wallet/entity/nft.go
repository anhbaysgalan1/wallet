package entity

import (
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NftData struct {
	GameName      string `json:"game_name,omitempty" bson:"game_name"`             // 游戏类型
	NftGameToken  string `json:"nft_game_token,omitempty" bson:"nft_game_token"`   // nft 游戏唯一标识
	NftBlockToken string `json:"nft_block_token,omitempty" bson:"nft_block_token"` // nft token唯一标识
	Level         uint64 `json:"level,omitempty" bson:"level"`                     // 等级
	Num           uint64 `json:"num,omitempty" bson:"num"`                         // 数量
}

func (nd NftData) IsEmpty() bool {
	return nd == NftData{}
}

type NftOwner struct {
	Id            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // 订单id
	GameId        string             `json:"game_id,omitempty" bson:"game_id"`               // 游戏id
	NftToken      string             `json:"nft_token,omitempty" bson:"nft_token"`           // nft唯一标识
	GameToken     string             `json:"game_token,omitempty" bson:"game_token"`         // nft在游戏唯一标识(所有修改，创建以游戏唯一标识为准)
	CreateHash    string             `json:"create_hash,omitempty" bson:"create_hash"`       // 创建hash(用于判断)
	ContractToken string             `json:"contract_token,omitempty" bson:"contract_token"` // 合约标识
	OwnerAddress  string             `json:"owner_address,omitempty" bson:"owner_address"`   // 拥有者地址
	CreateTime    int64              `json:"create_time,omitempty" bson:"create_time"`       // 创建时间
	UpdateTime    int64              `json:"update_time,omitempty" bson:"update_time"`       // 修改时间
	NftData       NftData            `json:"nft_data,omitempty" bson:"nft_data"`             // nft data
	Uid           uint64             `json:"uid,omitempty" bson:"-"`                         // 拥有者id
	Status        NftOwner_Status    `json:"status,omitempty" bson:"status"`                 // 拥有者状态
}

type NftOwner_Status uint16 // 拥有者状态

const (
	NftOwnerStatus_Available NftOwner_Status = 1 // 可以被操作
	NftOwnerStatus_Cash      NftOwner_Status = 2 // 提现中
)

func (b NftOwner) ToPb() *walletPb.NftInfo {
	return &walletPb.NftInfo{
		Uid:           b.Uid,
		OwnerAddress:  b.OwnerAddress,
		GameId:        b.GameId,
		NftGameToken:  b.GameToken,
		NftChainToken: b.NftToken,
		ContractToken: b.ContractToken,
		Level:         b.NftData.Level,
	}
}

type NftContract struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id"`                                    // 订单id
	GameId             string             `json:"game_id,omitempty" bson:"game_id"`                           // 游戏标识
	ContractToken      string             `json:"contract_token,omitempty" bson:"contract_token"`             // 合约信息
	NftToken           string             `json:"nft_token,omitempty" bson:"nft_token"`                       // nft唯一标识
	GameToken          string             `json:"game_token,omitempty" bson:"game_token"`                     // nft游戏内唯一标识
	OwnerAddress       string             `json:"owner_address,omitempty" bson:"owner_address"`               // 拥有者地址
	ReallyOwnerAddress string             `json:"really_owner_address,omitempty" bson:"really_owner_address"` // 真实拥有者地址
	CreateTime         int64              `json:"create_time,omitempty" bson:"create_time"`                   // 创建时间
	UpdateTime         int64              `json:"update_time,omitempty" bson:"update_time"`                   // 修改时间
	DeleteTime         int64              `json:"delete_time,omitempty" bson:"delete_time"`                   // 删除时间
	Uid                uint64             `json:"uid,omitempty" bson:"-"`                                     // 拥有者id
}

func (b NftContract) ToPb() *walletPb.NftInfo {
	return &walletPb.NftInfo{
		Uid:           b.Uid,
		OwnerAddress:  b.OwnerAddress,
		GameId:        b.GameId,
		NftGameToken:  b.GameToken,
		NftChainToken: b.NftToken,
		ContractToken: b.ContractToken,
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
