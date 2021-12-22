package service

import (
	"context"
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/big"
	"strings"
	"tp_wallet/config"
	"tp_wallet/internal/block_chain/chain/transfer"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

func (srv WalletSrv) NftCreate(ctx context.Context, info *walletPb.NftInfo) (*walletPb.Hash, error) {
	info.ContractAddress = strings.ToLower(info.ContractAddress)
	info.OwnerAddress = strings.ToLower(info.OwnerAddress)
	var err error
	var hash string
	var nonce uint64
	var billType entity.BillType
	if _, ok := config.SysAccountMap[info.ContractToken]; !ok {
		return nil, hcode.ErrContractNftAttributionForAddr
	}
	var toUid uint64
	var toAddr string
	if info.Uid == 0 {
		toUid = config.SysAccountMap[info.ContractToken].SysUid
		toAddr = config.SysAccountMap[info.ContractToken].AddrExpenditure[0]
	} else {
		// 校验地址是否属于uid
		var ownerUid uint64
		ownerUid, err = srv.Repo.UidGetByAddress(ctx, info.GetOwnerAddress())
		if err != nil {
			return nil, err
		}
		if ownerUid != info.GetUid() {
			log.GetLogger().Error("[NftCreate] uid different", zap.Any("req", info), zap.Uint64("owner uid", ownerUid))
			return nil, hcode.ErrPermissions
		}
		toUid = info.GetUid()
		toAddr = info.GetOwnerAddress()
	}
	// 获取nonce值
	nonce, err = srv.Repo.GetAndLockAddr(ctx, config.SysAccountMap[info.ContractToken].AddrCreate)
	if err != nil {
		return nil, err
	}
	req := transfer.InputForTransfer{
		FromAddress:     config.SysAccountMap[info.ContractToken].AddrCreate,
		ContractAddress: config.SysAccountMap[info.ContractToken].ContractAddress,
		Amount:          new(big.Int).SetUint64(0),
		ToAddress:       toAddr,
		GasLimit:        config.Fee.GasLimit,
		GasPrice:        config.Fee.GasPrice,
		Nonce:           nonce,
		Private:         "5799008c20d5a9bd55dd431fc102285321b216b52c7c77012e289b0d33111c18",
	}
	billType, err = config.GetBillTypeBySysUid(config.SysAccountMap[info.ContractToken].SysUid)
	if err != nil {
		log.GetLogger().Error("[NftCreate] config.GetBillTypeBySysUid failed", zap.Uint64("sys uid", config.SysAccountMap[info.ContractToken].SysUid))
		return nil, err
	}
	switch billType {
	case entity.BillType_Eip721:
		hash, err = srv.BlockChainSrv.RacingBoatCreateNft(ctx, req, info.NftGameToken, uint8(info.Level), config.BlockBusiness.KeyTransfer)
		if err != nil {
			log.GetLogger().Error("[NftCreate] BlockChainSrv.RacingBoatCreateNft failed", zap.Any("req", req),
				zap.Any("info", info), zap.Error(err))
			return nil, err
		}
	case entity.BillType_Eip1155:
		hash, err = srv.BlockChainSrv.MaterialCreate(ctx, req, new(big.Int).SetUint64(info.Num), info.NftGameToken, config.BlockBusiness.KeyTransfer)
		if err != nil {
			log.GetLogger().Error("[NftCreate] BlockChainSrv.RacingBoatCreateNft failed", zap.Any("req", req),
				zap.Any("info", info), zap.Error(err))
			return nil, err
		}
	default:
		log.GetLogger().Error("[NftCreate] config.GetBillTypeBySysUid failed", zap.Uint64("sys uid", config.SysAccountMap[info.ContractToken].SysUid))
		return nil, hcode.ErrBillType
	}
	if len(hash) > 0 {
		// 添加nonce
		_ = srv.Repo.NonceIncr(ctx, config.SysAccountMap[info.ContractToken].AddrCreate, 1)
		// 解锁
		srv.Repo.UnlockNonceAddr(ctx, config.SysAccountMap[info.ContractToken].AddrCreate)
	}
	newBill := &entity.Bill{
		Id:             primitive.NewObjectID(),
		NumericalOrder: uint64(tool.GetSnowFlake().GetId()),
		Uid:            toUid,
		TransferType:   walletPb.TransferType_NftCreate,
		BillStatus:     walletPb.BillStatus_Pending,
		FromUid:        config.SysAccountMap[info.ContractToken].SysUid,
		FromAddr:       config.SysAccountMap[info.ContractToken].AddrCreate,
		ToUid:          toUid,
		ToAddr:         toAddr,
		Hash:           hash,
		BillType:       billType,
		ContractRecord: entity.ContractRecord{
			ContractType: info.ContractToken,
			ContractAddr: config.SysAccountMap[info.ContractToken].ContractAddress,
			GameId:       info.GameId,
			NftToken:     "",
			GameToken:    info.GetNftGameToken(),
			Num:          info.GetNum(),
		},
	}
	// 创建nft owner
	var nftData = entity.NftData{
		GameName:     info.GameId,
		NftGameToken: info.NftGameToken,
		Level:        info.Level,
		Num:          info.GetNum(),
	}
	err = srv.Repo.NftCreate(ctx, newBill, nftData)
	if err != nil {
		return nil, err
	}
	return &walletPb.Hash{Hash: hash}, err
}

// NftGetAttribution 获取拥有者属性
func (srv WalletSrv) NftGetAttribution(ctx context.Context, token *walletPb.NftToken) (*walletPb.NftInfo, error) {
	var owner *entity.NftOwner
	var err error
	owner, err = srv.Repo.NftOwnerGetByNftToken(ctx, token.GetContractToken(), token.GetNftToken())
	if err != nil {
		return nil, err
	}
	return owner.ToPb(), nil
}

func (srv WalletSrv) NftGetByUid(ctx context.Context, req *walletPb.NftGetByUidReq) (*walletPb.NftInfoS, error) {
	if req.Page == nil && req.Page.Limit == 0 {
		return nil, hcode.ErrParameter
	}
	// 获取玩家地址
	var owners []*entity.NftOwner
	var err error
	owners, err = srv.Repo.NftOwnerGetByAddr(ctx, req.GetAddr(), req.Page)
	if err != nil {
		return nil, err
	}
	var result = new(walletPb.NftInfoS)
	result.Nfts = make(map[string]*walletPb.NftInfo)
	for _, item := range owners {
		if _, ok := config.SysAccountMap[item.ContractToken]; ok {
			if nftType, _ := config.GetBillTypeBySysUid(config.SysAccountMap[item.ContractToken].SysUid); nftType == entity.BillType_Eip1155 {
				if item.NftData.Num == 0 { // 用户并没有持有nft
					continue
				}
			}
			result.Nfts[item.NftToken] = item.ToPb()
		}
	}
	return result, nil
}

func (srv WalletSrv) NftCash(ctx context.Context, req *walletPb.NftCashReq) (*walletPb.Empty, error) {
	req.ToAddr = strings.ToLower(req.ToAddr)
	if _, ok := config.SysAccountMap[req.GetContractType()]; !ok {
		log.GetLogger().Error("[NftCash] parameters failed", zap.Any("req", req))
		return nil, hcode.ErrContractUnsupported
	}
	// nft是否存在我们系统账号
	var owner *entity.NftOwner
	var err error
	owner, err = srv.Repo.NftOwnerGetByNftToken(ctx, req.GetContractType(), req.GetNftToken())
	if err != nil {
		return nil, err
	}
	if owner.OwnerAddress != config.SysAccountMap[req.GetContractType()].AddrExpenditure[0] {
		log.GetLogger().Error("[NftCash] owner address is not sys address", zap.Any("owner", owner), zap.Any("req", req))
		return nil, hcode.ErrPermissions
	}
	// nft目前的状态
	if owner.Status != entity.NftOwnerStatus_Available {
		log.GetLogger().Error("[NftCash] entity NftOwnerStatus_Available failed", zap.Any("owner", owner), zap.Any("req", req))
		return nil, hcode.ErrPermissions
	}
	newBill := &entity.Bill{
		Uid:    req.Uid,
		ToUid:  req.Uid,
		ToAddr: req.GetToAddr(),
		ContractRecord: entity.ContractRecord{
			ContractType: req.ContractType,
			ContractAddr: config.SysAccountMap[req.GetContractType()].ContractAddress,
			NftToken:     req.NftToken,
			GameToken:    owner.GameToken,
			Num:          req.GetNum(),
		},
	}
	newBill.BillType, err = config.GetBillTypeBySysUid(config.SysAccountMap[req.ContractType].SysUid)
	if err != nil {
		log.GetLogger().Error("[NftCreate] config.GetBillTypeBySysUid failed", zap.Uint64("sys uid", config.SysAccountMap[req.ContractType].SysUid))
		return nil, err
	}
	return srv.Repo.TransferNftCash(ctx, newBill)
}

func (srv WalletSrv) GetSysContractAddr(ctx context.Context, req *walletPb.ContractType) (*walletPb.AddrResp, error) {
	var result = new(walletPb.AddrResp)

	if _, ok := config.SysAccountMap[req.GetContractType()]; !ok {
		return nil, hcode.ErrSysAccountNotFound
	}
	result.Addr = config.SysAccountMap[req.GetContractType()].AddrIncome
	return result, nil
}
