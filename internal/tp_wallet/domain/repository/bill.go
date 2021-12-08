package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/big"
	"tp_wallet/config"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

// BillH2OGetForAsync 获取链上异步未结算账单
func (repo RepositoryStruct) BillH2OGetForAsync(ctx context.Context) ([]*entity.Bill, error) {
	return repo.Db.BillH2OGetForQueue(ctx)
}

// BillNftGetForAsync 获取链上异步未结算账单
func (repo RepositoryStruct) BillNftGetForAsync(ctx context.Context) ([]*entity.Bill, error) {
	return repo.Db.BillNftGetForQueue(ctx)
}

// BillDealWith 处理一笔订单
func (repo RepositoryStruct) BillDealWith(ctx context.Context, bill *entity.Bill) (err error) {
	if bill == nil {
		return hcode.ErrInternalParameter
	}
	// 在这里组装数据
	if len(bill.FromAddr) > 0 {
		bill.FromUid, _ = repo.UidGetByAddress(ctx, bill.FromAddr)
	}
	if len(bill.ToAddr) > 0 {
		bill.ToUid, _ = repo.UidGetByAddress(ctx, bill.ToAddr)
	}
	switch bill.BillStatus {
	case dto.BillStatus_Pending:
		return repo.BillDealWithPending(ctx, bill)
	case dto.BillStatus_Success:
		return repo.BillDealWithSuccess(ctx, bill)
	case dto.BillStatus_Failed:
		return repo.BillDealWithFailed(ctx, bill)
	default:
		return hcode.ErrInternalParameter
	}
}

// BillDealWithPending 处理上链中订单
func (repo RepositoryStruct) BillDealWithPending(ctx context.Context, bill *entity.Bill) error {
	switch bill.BillType {
	case entity.BillType_Currency:
		return repo.BillCurrencyDealWithPending(ctx, bill)
	case entity.BillType_Nft:
		return repo.BillNftDealWithPending(ctx, bill)
	default:
		return hcode.ErrBillType
	}
}

func (repo RepositoryStruct) BillCurrencyDealWithPending(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || bill.NumericalOrder == 0 || bill.Id.IsZero() {
		log.GetLogger().Error("[BillH2ODealWithPending] hash length or id", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 修改订单hash以及订单状态
	bill.DealWithBillPending()
	var filter = map[string]interface{}{"numerical_order": bill.NumericalOrder}
	err = repo.Db.BillSet(ctx, bill, filter)
	return err
}

func (repo RepositoryStruct) BillNftDealWithPending(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.FromAddr) == 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillNftDealWithPending] hash length or id", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	bill.DealWithBillPending()
	var filter = map[string]interface{}{"numerical_order": bill.NumericalOrder}
	err = repo.Db.BillSet(ctx, bill, filter)
	return err
}

func (repo RepositoryStruct) BillDealWithSuccess(ctx context.Context, bill *entity.Bill) error {
	if len(bill.Hash) == 0 || len(bill.ToAddr) == 0 || len(bill.FromAddr) == 0 {
		return hcode.ErrInternalParameter
	}
	var billForStore *entity.Bill
	billForStore, _ = repo.Db.BillGetByHash(ctx, bill.Hash)
	// 获取合约类型
	if _, ok := config.BlockSysAddrToUid[bill.ToAddr]; ok {
		bill.ContractType = config.SysAccountMapByUid[config.BlockSysAddrToUid[bill.ToAddr]].SysToken
	}
	if _, ok := config.BlockSysAddrToUid[bill.FromAddr]; ok {
		bill.ContractType = config.SysAccountMapByUid[config.BlockSysAddrToUid[bill.FromAddr]].SysToken
	}
	if billForStore != nil {
		if billForStore.BillStatus != dto.BillStatus_Pending || billForStore.FromAddr != bill.FromAddr || billForStore.ToAddr != bill.ToAddr {
			log.GetLogger().Error("[BillDealWithSuccess] bill status check failed",
				zap.Any("inset bill", bill),
				zap.Any("store bill", billForStore))
			return nil
		}
		// 在这里判断订单具体类型，因为区块链传过来的类型，提现会被划分为转账
		bill.Id = billForStore.Id
		bill.BillType = billForStore.BillType
	}
	switch bill.BillType {
	case dto.TransferType_H2OTransfer:
		log.GetLogger().Info("[BillDealWithSuccess] exterior bill", zap.Any("bill", bill))
		return nil
	case dto.TransferType_H2OCASH:
		return repo.BillDealWithH2OCash(ctx, bill)
	case dto.TransferType_H2OCHARGE:
		return repo.BillDealWithH2OCharge(ctx, bill)
	case dto.TransferType_NftTransfer:
		// todo: 需要修改nft拥有者 & 删除授权给我们的nft道具信息
		panic(hcode.ErrServer)
	case dto.TransferType_NftApproval:
		// todo: 授权给我们，需要展示他授权给我们的内容
		panic(hcode.ErrServer)
	case dto.TransferType_NftApprovalTransfer:
		// todo: 需要修改nft拥有者 & 删除授权给我们的nft道具信息
		panic(hcode.ErrServer)
	case dto.TransferType_NftCreate:
		return repo.BillDealWithNftCreate(ctx, bill)
	case dto.TransferType_NftCASH:
		return repo.BillDealWithNftCash(ctx, bill)
	case dto.TransferType_NftCHARGE:
		return repo.BillDealWithNftCharge(ctx, bill)
	default:
		log.GetLogger().Error("[BillDealWithSuccess] bill type check failed",
			zap.Any("inset bill", bill),
			zap.Any("store bill", billForStore))
		return hcode.ErrInternalParameter
	}
}

// BillDealWithH2OCash 处理提现
func (repo RepositoryStruct) BillDealWithH2OCash(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillDealWithH2OCash] hash length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 修改系统用户金额
	bill.FromBeforeBalance, bill.FromAfterBalance, err = repo.BalanceSet(ctx, bill.From, bill.Amount, false)
	// 修改账单状态
	err = repo.Db.BillSetByHash(ctx, &entity.Bill{
		Hash:              bill.Hash,
		BillStatus:        dto.BillStatus_Success,
		ReceivedAmount:    bill.ReceivedAmount,
		FromBeforeBalance: bill.FromBeforeBalance,
		FromAfterBalance:  bill.FromAfterBalance,
	})
	if err != nil {
		return err
	}
	return err
}

// BillDealWithH2OCharge 处理充值订单
func (repo RepositoryStruct) BillDealWithH2OCharge(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) == 0 || len(bill.ToAddr) == 0 || len(bill.FromAddr) == 0 || len(bill.ReceivedAmount) == 0 {
		log.GetLogger().Error("[BillDealWithH2OCharge] hash, amount or from addr length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// To用户id
	bill.To, err = repo.UidGetByAddress(ctx, bill.ToAddr)
	if err != nil {
		return err
	}
	bill.From, _ = repo.UidGetByAddress(ctx, bill.FromAddr)
	bill.BillType = dto.TransferType_H2OCHARGE
	bill.BillStatus = dto.BillStatus_Success
	bill.Amount = bill.ReceivedAmount
	bill.IsBalanceTrade = true
	// 修改用户金额
	bill.ToBeforeBalance, bill.ToAfterBalance, err = repo.BalanceSet(ctx, bill.To, bill.Amount, true)
	// 创建订单
	err = repo.Db.BillCreate(ctx, bill)
	return err
}

// BillDealWithNftCreate 处理创建
func (repo RepositoryStruct) BillDealWithNftCreate(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillDealWithNftCash] hash length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 获取nft拥有者
	var nftOwner *entity.NftOwner
	nftOwner, err = repo.Db.NftOwnerGetByCreateHash(ctx, bill.Hash)
	if err != nil {
		return err
	}
	nftOwner.NftToken = bill.NftToken
	nftOwner.NftData.NftBlockToken = bill.NftToken
	// 获取nft合约绑定
	var nftContract *entity.NftContract
	nftContract, err = repo.Db.NftContractGetByGameToken(ctx, bill.ContractType, nftOwner.GameToken)
	if err != nil {
		return err
	}
	nftContract.NftToken = bill.NftToken
	nftContract.NftData.NftBlockToken = bill.NftToken
	// 修改账单状态
	err = repo.Db.BillSetByHash(ctx, &entity.Bill{
		Hash:       bill.Hash,
		BillStatus: dto.BillStatus_Success,
		NftToken:   bill.NftToken,
	})
	if err != nil {
		return err
	}
	// 修改nft拥有者
	err = repo.Db.NftOwnerSetByGameToken(ctx, nftOwner)
	if err != nil {
		return err
	}
	err = repo.Db.NftContractSetByGameToken(ctx, nftContract)
	if err != nil {
		return err
	}
	return err
}

// BillDealWithNftCash 处理提现
func (repo RepositoryStruct) BillDealWithNftCash(ctx context.Context, bill *entity.Bill) error {
	var err error
	if len(bill.Hash) <= 0 || len(bill.ToAddr) == 0 {
		log.GetLogger().Error("[BillDealWithNftCash] hash length", zap.Any("bill", bill))
		return hcode.ErrInternalParameter
	}
	// 修改账单状态
	err = repo.Db.BillSetByHash(ctx, &entity.Bill{
		Hash:       bill.Hash,
		BillStatus: dto.BillStatus_Success,
		NftToken:   bill.NftToken,
	})
	if err != nil {
		return err
	}
	// 修改nft拥有者
	var now = tool.GetTimeUnixMilli()
	_ = repo.Db.NftOwnerSetByToken(ctx, &entity.NftOwner{
		NftToken:     bill.NftToken,
		OwnerAddress: bill.ToAddr,
		CreateTime:   now,
		UpdateTime:   now,
	})
	return err
}

// BillDealWithNftCharge 处理充值订单
func (repo RepositoryStruct) BillDealWithNftCharge(ctx context.Context, bill *entity.Bill) error {
	var err error
	// 判断nft是否重复
	var sysUid uint64
	sysUid, err = repo.UidGetByAddress(ctx, bill.ToAddr)
	if err != nil {
		return err
	}
	var storeNftToken *entity.NftContract
	var storeNftOwner *entity.NftOwner
	storeNftToken, err = repo.Db.NftContractGetByNftToken(ctx, config.SysAccountMapByUid[sysUid].SysToken, bill.NftToken)
	if err != nil && !errors.Is(err, hcode.ErrNftNotFound) {
		return err
	}
	if storeNftToken != nil {
		if storeNftToken.DeleteTime <= 0 {
			log.GetLogger().Error("[BillDealWithNftCharge] nft already exists", zap.Any("bill", bill), zap.Any("nft", bill))
			return hcode.ErrServer
		}
	}
	storeNftOwner, err = repo.Db.NftOwnerGetByNftToken(ctx, config.SysAccountMapByUid[sysUid].SysToken, bill.NftToken)
	if err != nil && !errors.Is(err, hcode.ErrNftNotFound) {
		return err
	}
	// 创建订单
	err = repo.Db.BillCreate(ctx, bill)
	// 修改nft拥有者
	var now = tool.GetTimeUnixMilli()
	// 修改nft拥有者
	if err = repo.Db.NftOwnerSetByToken(ctx, &entity.NftOwner{
		NftToken:      bill.NftToken,
		ContractToken: config.SysAccountMapByUid[sysUid].SysToken,
		OwnerAddress:  bill.ToAddr,
		UpdateTime:    now,
	}); err != nil {
		return err
	}
	// 充值到合约
	if err = repo.Db.NftContractCreate(ctx, &entity.NftContract{
		Id:                 primitive.NewObjectID(),
		NftToken:           bill.NftToken,
		GameToken:          storeNftOwner.NftData.NftGameToken,
		ContractToken:      bill.ContractType,
		OwnerAddress:       bill.FromAddr,
		ReallyOwnerAddress: bill.ToAddr,
		CreateTime:         now,
		UpdateTime:         now,
		NftData:            storeNftOwner.NftData,
	}); err != nil {
		return err
	}
	return err
}

// BillDealWithFailed 处理上链失败订单
func (repo RepositoryStruct) BillDealWithFailed(ctx context.Context, bill *entity.Bill) error {
	var err error
	// 获取账单
	bill, err = repo.Db.BillGetById(ctx, bill.Id)
	if err != nil {
		return err
	}
	var amount = new(big.Int)
	amount.SetString(bill.Amount, 10)
	var toBeforeBalance, toAfterBalance string
	if amount.Cmp(big.NewInt(0)) > 0 {
		if bill.IsBalanceTrade {
			// 给to退钱 因为from是系统账户
			toBeforeBalance, toAfterBalance, err = repo.BalanceSet(ctx, bill.To, bill.Amount, true)
		}
	}
	var remarks = entity.ErrorBillRemark{ // 存储错误信息
		To:              bill.To,
		ToBeforeBalance: toBeforeBalance,
		ToAfterBalance:  toAfterBalance,
		Data:            bill.Remark,
	}.ToJson()
	// 修改订单状态
	err = repo.Db.BillSetById(ctx, &entity.Bill{Id: bill.Id, BillStatus: dto.BillStatus_Failed, Remark: remarks})
	return err
}
