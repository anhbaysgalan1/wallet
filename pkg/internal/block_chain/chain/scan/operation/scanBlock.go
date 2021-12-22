package operation

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"tp_wallet/internal/block_chain/chain/bsc"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/block_chain/chain/scan/orm"
	"tp_wallet/pkg/log"
)

func ScanBlockForBsc() (string, bool, error) {
	//查看mongo里面有没有数据
	ma, err := orm.GetHeightByMongo(orm.CollectionBsc)
	if err != nil {
		return "", false, err
	}
	ma.Height += 1

	bi, err := bsc.GetBlockInfo(ma.Height, common.ScanKey)
	if err != nil {
		return "", false, err
	}

	if bi.Result.Hash == "" || bi.Result.Number == "" {
		return "", false, errors.New(common.GetBlockErrNotFound)
	}

	hi, err := strconv.ParseUint(bi.Result.Number, 0, 0)
	if err != nil {
		log.GetLogger().Error("strconv.ParseUint(bi.Result.Number)", zap.Error(err), zap.Any("blockNumber", bi.Result.Number))
		return "", false, errors.New(common.GetBlockErrNotFound)
	}

	if hi != ma.Height {
		return "", false, errors.New(common.GetBlockErrNotFound)
	}

	if len(bi.Result.Transactions) != 0 {
		for _, v := range bi.Result.Transactions {
			//如果交易获取失败，重新查询该块，防止丢失数据
			errC := checkTransactionForBsc(strings.ToLower(v.From), strings.ToLower(v.To), v.Hash, v.Input, v.Nonce, v.Value, v.BlockNumber)
			if errC != nil {
				return "", false, errC
			}
		}
	}
	return bi.Result.Number, true, nil
}

// checkTransactionForBsc 判断to是不是相关合约地址
func checkTransactionForBsc(from, to, txh, input, nonce, value, blockNumber string) error {

	if from == common.BNBWithdrawAddress &&
		to != strings.ToLower(common.RacingBoatContractAddress) &&
		to != strings.ToLower(common.RacerContractAddress) &&
		to != strings.ToLower(common.FFCoinContractAddress) &&
		to != strings.ToLower(common.F1CoinContractAddress) {
		return PushBNBTransaction(from, to, value, txh, nonce, blockNumber, common.BNBWithdrawCode)
	}

	switch to {
	//case common.H2OContractAddress:
	//	addr1, addr2, amount, code, err := CheckH2OInput(input)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if code == common.UnknownCode {
	//		return nil
	//	}
	//
	//	if code == common.H2OTransferFromCode {
	//		err = PushH2oTransaction(addr1, addr2, to, amount, txh, nonce, blockNumber, code)
	//	} else {
	//		err = PushH2oTransaction(from, addr1, to, amount, txh, nonce, blockNumber, code)
	//	}
	//	return err
	case strings.ToLower(common.RacingBoatContractAddress):
		addr1, addr2, NftID, code, err := CheckRacingBoatInput(input)
		if err != nil {
			return nil
		}

		if code == common.UnknownCode {
			return nil
		}

		if code == common.RowingNftTransferFromCode {
			err = PushRowingTransaction(addr1, addr2, to, NftID, txh, nonce, blockNumber, code)
		} else if code == common.RowingNftCreateCode {
			err = PushCreateRowingNft(from, to, txh, nonce, blockNumber, code)
		} else {
			err = PushRowingTransaction(from, addr1, to, NftID, txh, nonce, blockNumber, code)
		}
		return err
	case strings.ToLower(common.RacerContractAddress):
		addr1, addr2, NftID, code, err := CheckRacerInput(input)
		if err != nil {
			return nil
		}

		if code == common.UnknownCode {
			return nil
		}

		if code == common.RacerNftTransferFromCode {
			err = PushRacerTransaction(addr1, addr2, to, NftID, txh, nonce, blockNumber, code)
		} else if code == common.RacerNftCreateCode {
			err = PushCreateRacerNft(from, to, txh, nonce, blockNumber, code)
		} else {
			err = PushRacerTransaction(from, addr1, to, NftID, txh, nonce, blockNumber, code)
		}
		return err
	case strings.ToLower(common.FFCoinContractAddress):
		addr1, addr2, amount, code, err := CheckFFInput(input)
		if err != nil {
			return nil
		}

		if code == common.UnknownCode {
			return nil
		}

		if code == common.FFTransferFromCode {
			err = PushFFTransaction(addr1, addr2, to, amount, txh, nonce, blockNumber, code)
		} else {
			err = PushFFTransaction(from, addr1, to, amount, txh, nonce, blockNumber, code)
		}
		return err
	case strings.ToLower(common.F1CoinContractAddress):
		addr1, addr2, amount, code, err := CheckF1Input(input)
		if err != nil {
			return nil
		}

		if code == common.UnknownCode {
			return nil
		}

		if code == common.F1TransferFromCode {
			err = PushF1Transaction(addr1, addr2, to, amount, txh, nonce, blockNumber, code)
		} else {
			err = PushF1Transaction(from, addr1, to, amount, txh, nonce, blockNumber, code)
		}
		return err
	case strings.ToLower(common.MaterialContractAddress):
		_, _to, _nftID, _amount, code, err := CheckMaterialInput(input)
		if err != nil {
			return nil
		}

		if code == common.UnknownCode {
			return nil
		}

		if code == common.MaterialNftCreateCode {
			err = PushCreateMaterialNft(from, "", txh, nonce, blockNumber, common.MaterialNftCreateCode)
		} else if code == common.MaterialNftExpandCode {
			err = PushMaterialTransaction(
				from, _to,
				strings.ToLower(common.MaterialContractAddress),
				_nftID,
				_amount,
				txh, nonce, blockNumber, common.MaterialNftExpandCode)
		} else if code == common.MaterialNftTransferCode {
			err = PushMaterialTransaction(
				from,
				_to,
				strings.ToLower(common.MaterialContractAddress),
				_nftID,
				_amount,
				txh, nonce, blockNumber, common.MaterialNftTransferCode)
		}
		return err
	case strings.ToLower(common.BNBRechargeAddress):
		return PushBNBTransaction(from, to, value, txh, nonce, blockNumber, common.BNBRechargeCode)
	default:
		return nil
	}
}
