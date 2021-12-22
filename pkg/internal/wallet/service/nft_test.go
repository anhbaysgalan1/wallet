package service

import (
	"fmt"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"testing"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/wallet/entity"
	"tp_wallet/pkg/tool"
)

func TestNftCreate(t *testing.T) {
	result, err := walletSrv.NftCreate(ctx, &walletPb.NftInfo{
		Uid:             0,
		OwnerAddress:    "",
		GameId:          "1",
		NftGameToken:    "test_game_token_" + tool.Int64ToStr(tool.GetTimeUnixMilli()),
		NftChainToken:   "",
		ContractAddress: "",
		ContractToken:   "Rowing",
		Level:           5,
		Num:             0,
	})
	if err != nil {
		fmt.Println("failed ======>", err)
	}
	fmt.Println("success =======>", result)
}

func TestNftCreateFoBill(t *testing.T) {
	var blockData = common.Nft721CreateForPush{
		From:        "0x00000000000001",
		To:          "0x7cE8102020f45e1451ba93B7D7997AD6aF0ED56b",
		Contract:    "0x00000000000000",
		Nonce:       "",
		Hash:        "test_hash1639534999379",
		NftToken:    "test_nft_token" + tool.Int64ToStr(tool.GetTimeUnixMilli()),
		PropsName:   "",
		StarRating:  "",
		Status:      "success",
		BlockNumber: "",
	}
	var newBill = &walletPb.BillInfo{
		BillType:     int64(entity.BillType_Eip721),
		TransferType: walletPb.TransferType_NftCHARGE,
		Hash:         blockData.Hash,
		FromAddr:     blockData.From,
		ToAddr:       blockData.To,
		ContractRecord: &walletPb.ContractRecord{
			ContractAddr: blockData.Contract,
			NftToken:     blockData.NftToken,
		},
	}
	if blockData.Status == "success" {
		newBill.BillStatus = walletPb.BillStatus_Success
	} else {
		newBill.BillStatus = walletPb.BillStatus_Failed
	}
	result, err := walletSrv.DealWithBill(ctx, newBill)
	if err != nil {
		fmt.Println("failed =====>", err)
	}
	fmt.Println("success ========>", result)
}

func Test_NftCash(t *testing.T) {
	result, err := walletSrv.NftCash(ctx, &walletPb.NftCashReq{
		Uid:          2502500,
		ToAddr:       "0xed2dfc20e4647381340ecb0f5d3f8e452fcb97f7",
		NftToken:     "6",
		ContractType: "Rowing",
	})
	if err != nil {
		fmt.Println("failed =======>", err)
	}
	fmt.Println("success =======>", result)
}

func Test_NftCashToBill(t *testing.T) {
	var blockData = common.Nft721TransactionForPush{
		From:        "0x7cE8102020f45e1451ba93B7D7997AD6aF0ED56b",
		To:          "0x0000000000000000000000002000000",
		Contract:    "0x00000000000000",
		Nonce:       "",
		Hash:        "test_hash1639550045023",
		NftToken:    "test_nft_token1639536673669",
		Status:      "success",
		BlockNumber: "",
	}
	var newBill = &walletPb.BillInfo{
		BillType:     int64(entity.BillType_Eip721),
		TransferType: walletPb.TransferType_NftCHARGE,
		Hash:         blockData.Hash,
		FromAddr:     blockData.From,
		ToAddr:       blockData.To,
		ContractRecord: &walletPb.ContractRecord{
			ContractAddr: blockData.Contract,
			NftToken:     blockData.NftToken,
		},
	}
	if blockData.Status == "success" {
		newBill.BillStatus = walletPb.BillStatus_Success
	} else {
		newBill.BillStatus = walletPb.BillStatus_Failed
	}
	result, err := walletSrv.DealWithBill(ctx, newBill)
	if err != nil {
		fmt.Println("failed =====>", err)
	}
	fmt.Println("success ========>", result)
}

func Test_NftChargeToBill(t *testing.T) {
	var blockData = common.Nft721TransactionForPush{
		From:        "0x0000000000000000000000002000000",
		To:          "0x7cE8102020f45e1451ba93B7D7997AD6aF0ED56b",
		Contract:    "0x00000000000000",
		Nonce:       "",
		Hash:        "test_hash" + tool.Int64ToStr(tool.GetTimeUnixMilli()),
		NftToken:    "test_nft_token1639536673669",
		Status:      "success",
		BlockNumber: "",
	}
	var newBill = &walletPb.BillInfo{
		BillType:     int64(entity.BillType_Eip721),
		TransferType: walletPb.TransferType_NftCHARGE,
		Hash:         blockData.Hash,
		FromAddr:     blockData.From,
		ToAddr:       blockData.To,
		ContractRecord: &walletPb.ContractRecord{
			ContractAddr: blockData.Contract,
			NftToken:     blockData.NftToken,
		},
	}
	if blockData.Status == "success" {
		newBill.BillStatus = walletPb.BillStatus_Success
	} else {
		newBill.BillStatus = walletPb.BillStatus_Failed
	}
	result, err := walletSrv.DealWithBill(ctx, newBill)
	if err != nil {
		fmt.Println("failed =====>", err)
	}
	fmt.Println("success ========>", result)
}

func Test_NftGetByUid(t *testing.T) {
	result, err := walletSrv.NftGetByUid(ctx, &walletPb.NftGetByUidReq{
		Uid:  2502500,
		Addr: "0xed2dfc20e4647381340ecb0f5d3f8e452fcb97f7",
		Page: &walletPb.Page{
			Limit:  10,
			Offset: 0,
		},
	})
	if err != nil {
		fmt.Println("failed ======>", err)
	}
	fmt.Println("success =======>", result)
}
