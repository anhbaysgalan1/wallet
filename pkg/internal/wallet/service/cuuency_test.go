package service

import (
	"fmt"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.uber.org/zap"
	"testing"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

func TestWalletSrv_AccountGet(t *testing.T) {
	//for i := 2000000; i <= 2000500; i++ {
	i := 2502500
	var uidReq = &walletPb.UidReq{
		Uid:      uint64(i),
		Addr:     "0xeD2dFc20E4647381340Ecb0F5D3F8e452fCb97f7",
		Currency: "ff",
	}
	if result, err := walletSrv.AccountGet(ctx, uidReq); err != nil {
		log.GetLogger().Error("[ConsumeClaim] failed", zap.Any("uidReq", uidReq), zap.Error(err))
	} else {
		fmt.Println("success ====>", result)
	}
	//}
}

func TestWalletSrv_BalanceGet(t *testing.T) {
	var uidReq = &walletPb.UidReq{
		Uid:  2000000,
		Addr: "0x0000000000000000000000002",
	}
	if result, err := walletSrv.BalanceGet(ctx, uidReq); err != nil {
		log.GetLogger().Error("[ConsumeClaim] failed", zap.Any("uidReq", uidReq), zap.Error(err))
	} else {
		fmt.Println("success ====>", result)
	}
}

func TestWalletSrv_TransferCurrencyForOffline(t *testing.T) {
	_, err := walletSrv.TransferCurrencyForOffline(ctx, &walletPb.TransferForOfflineReq{
		From:         1,
		To:           2000000,
		Cid:          1,
		Currency:     "usdt",
		Amount:       "1000",
		TransferType: walletPb.TransferType_CurrencyTransfer,
	})
	if err != nil {
		fmt.Println("error ======>", err)
	} else {
		fmt.Println("success ======>")
	}
}

func TestGetSysTransferAddr(t *testing.T) {
	result, err := walletSrv.GetSysTransferAddr(ctx, &walletPb.CurrencyReq{Currency: "bnb"})
	if err != nil {
		fmt.Println("failed ====>", err)
	}
	fmt.Println("success =======>", result)
}

func TestWalletSrv_TransferCurrencyCash(t *testing.T) {
	result, err := walletSrv.TransferCurrencyCash(ctx, &walletPb.TransferCashReq{
		Uid:      2502500,
		ToAddr:   "0xed2dfc20e4647381340ecb0f5d3f8e452fcb97f7",
		Currency: "f1",
		Amount:   "12000000000000000000",
	})
	if err != nil {
		fmt.Println("failed ======>", err)
	} else {
		fmt.Println("success =======>", result)
	}
}

func TestWalletSrv_TransferCurrencyCashToBill(t *testing.T) {
	var blockData = common.H20TransactionForPush{
		From:        "0x_usdt_expenditure_1",
		To:          "0x0000000000000000000000002000000",
		Amount:      "98",
		Nonce:       "1",
		Hash:        "test_hash1639445045004",
		Status:      "success",
		BlockNumber: "1",
		Currency:    "usdt",
	}
	var newBill = &walletPb.BillInfo{
		BillType: int64(entity.BillType_Eip20),
		Hash:     blockData.Hash,
		FromAddr: blockData.From,
		ToAddr:   blockData.To,
		BalanceRecord: &walletPb.BalanceRecord{
			Amount:        "",
			ReceiveAmount: blockData.Amount,
			BeforeBalance: "",
			AfterBalance:  "",
			Currency:      blockData.Currency,
		},
	}
	if blockData.Status == "success" {
		newBill.BillStatus = walletPb.BillStatus_Success
	} else {
		newBill.BillStatus = walletPb.BillStatus_Failed
	}
	result, err := walletSrv.DealWithBill(ctx, newBill)
	if err != nil {
		fmt.Println("failed ======>", err)
	} else {
		fmt.Println("success =======>", result)
	}
}

func TestWalletSrv_TransferCurrencyChargeToBill(t *testing.T) {
	var blockData = common.H20TransactionForPush{
		From:        "0x0000000000000000000000002000000",
		To:          "0x_usdt_income",
		Amount:      "1000",
		Nonce:       "1",
		Hash:        "0x_test_hash_" + tool.Int64ToStr(tool.GetTimeUnixMilli()),
		Status:      "success",
		BlockNumber: "1",
		Currency:    "usdt",
	}
	var newBill = &walletPb.BillInfo{
		BillType: int64(entity.BillType_Eip20),
		Hash:     blockData.Hash,
		FromAddr: blockData.From,
		ToAddr:   blockData.To,
		BalanceRecord: &walletPb.BalanceRecord{
			Amount:        "",
			ReceiveAmount: blockData.Amount,
			BeforeBalance: "",
			AfterBalance:  "",
		},
	}
	if blockData.Status == "success" {
		newBill.BillStatus = walletPb.BillStatus_Success
	} else {
		newBill.BillStatus = walletPb.BillStatus_Failed
	}
	result, err := walletSrv.DealWithBill(ctx, newBill)
	if err != nil {
		fmt.Println("failed ======>", err)
	} else {
		fmt.Println("success =======>", result)
	}
}
