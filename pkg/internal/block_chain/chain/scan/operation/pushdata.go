package operation

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.uber.org/zap"
	"tp_wallet/internal/block_chain/chain/bsc"
	"tp_wallet/internal/block_chain/chain/contract/erc1155"
	"tp_wallet/internal/block_chain/chain/contract/erc721"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/block_chain/chain/scan/orm"
	"tp_wallet/pkg/log"
)

// PushRowingTransaction 推送
func PushRowingTransaction(from, to, contract, nftID, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.Nft721TransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}

	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.NftToken = nftID
	hi.Nonce = nonce
	hi.Hash = txh

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushRowingTransaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushRowingTransaction", zap.Any("hash", txh))

	err = orm.SetRowingNftByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set RowingTransaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

// PushRacerTransaction 推送
func PushRacerTransaction(from, to, contract, nftID, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.Nft721TransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}

	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.NftToken = nftID
	hi.Nonce = nonce
	hi.Hash = txh

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushRacerTransaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushRacerTransaction", zap.Any("hash", txh))

	err = orm.SetRacerNftByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set RacerTransaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

func PushCreateRacerNft(from, contract, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}

	var hi common.Nft721CreateForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.Nonce = nonce
	hi.Hash = txh
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	if len(txInfo.Result.Logs) > 0 {
		data, err := hexutil.Decode(txInfo.Result.Logs[0].Data)
		if err != nil {
			log.GetLogger().Error("hexutil.Decode(txInfo.Result.Logs[0].Data)", zap.Error(err), zap.Any("hash", txh))
		} else {
			announcer, tokenID, pn, sr, err := erc721.UnPackCreateAssetEven(data)
			if err != nil {
				log.GetLogger().Error("erc721.UnPackCreateAssetEven", zap.Error(err), zap.Any("hash", txh))
			} else {
				hi.To = announcer
				hi.NftToken = tokenID
				hi.StarRating = sr
				hi.PropsName = pn
			}
		}
	}

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicCreateNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushCreateRacerNft", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushCreateRacerNft", zap.Any("hash", txh))

	err = orm.SetRacerCreateNftByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set CreateRacerNft", zap.Error(err), zap.Any("hash", txh))
	}

	return nil
}

func PushCreateRowingNft(from, contract, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}

	var hi common.Nft721CreateForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.Nonce = nonce
	hi.Hash = txh
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	if len(txInfo.Result.Logs) > 0 {
		data, err := hexutil.Decode(txInfo.Result.Logs[0].Data)
		if err != nil {
			log.GetLogger().Error("hexutil.Decode(txInfo.Result.Logs[0].Data)", zap.Error(err), zap.Any("hash", txh))
		} else {
			announcer, tokenID, pn, sr, err := erc721.UnPackCreateAssetEven(data)
			if err != nil {
				log.GetLogger().Error("erc721.UnPackCreateAssetEven", zap.Error(err), zap.Any("hash", txh))
			} else {
				hi.To = announcer
				hi.NftToken = tokenID
				hi.StarRating = sr
				hi.PropsName = pn
			}
		}
	}

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicCreateNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushCreateRowingNft", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushCreateRowingNft", zap.Any("hash", txh))

	err = orm.SetRowingCreateNftByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set CreateRowingNft", zap.Error(err), zap.Any("hash", txh))
	}

	return nil
}

func PushCreateMaterialNft(from, contract, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}

	var hi common.Nft1155CreateForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.Nonce = nonce
	hi.Hash = txh
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	if len(txInfo.Result.Logs) > 0 {
		data, err := hexutil.Decode(txInfo.Result.Logs[0].Data)
		if err != nil {
			log.GetLogger().Error("hexutil.Decode(txInfo.Result.Logs[0].Data)", zap.Error(err), zap.Any("hash", txh))
		} else {
			announcer, tokenID, name, amount, err := erc1155.UnPackCreateAssetEven(data)
			if err != nil {
				log.GetLogger().Error("erc1155.UnPackCreateAssetEven", zap.Error(err), zap.Any("hash", txh))
			} else {
				hi.To = announcer
				hi.NftToken = tokenID
				hi.MaterialName = name
				hi.Amount = amount
			}
		}
	}

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicCreateNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushMaterialCreateNft", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushMaterialCreateNft", zap.Any("hash", txh))

	err = orm.SetMaterialCreateByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set MaterialCreate", zap.Error(err), zap.Any("hash", txh))
	}

	return nil
}

func PushMaterialTransaction(from, to, contract, nftID, amount, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.Nft1155TransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}

	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.NftToken = nftID
	hi.Nonce = nonce
	hi.Hash = txh
	hi.Amount = amount

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicNftForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushMaterialTransaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushMaterialTransaction", zap.Any("hash", txh))

	err = orm.SetMaterialByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set MaterialTransaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

func PushFFTransaction(from, to, contract, amount, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.FFTransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.Amount = amount
	hi.Nonce = nonce
	hi.Hash = txh
	hi.Currency = common.FFCurrency

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicH20ForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushFFTransaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushFFTransaction", zap.Any("hash", txh))

	err = orm.SetFFByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set FFTransaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

func PushF1Transaction(from, to, contract, amount, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.F1TransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.Amount = amount
	hi.Nonce = nonce
	hi.Hash = txh
	hi.Currency = common.F1Currency

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicH20ForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushF1Transaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushF1Transaction", zap.Any("hash", txh))

	err = orm.SetF1ByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set F1Transaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

func PushBNBTransaction(from, to, amount, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.BNBTransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.To = to
	hi.From = from
	hi.BlockNumber = blockNumber
	hi.Amount = amount
	hi.Nonce = nonce
	hi.Hash = txh
	hi.Currency = common.BnBCurrency

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicH20ForScan, pushData)
	if err != nil {
		log.GetLogger().Error("PushBNBTransaction", zap.Error(err), zap.Any("hash", txh))
		return err
	}
	log.GetLogger().Info("PushBNBTransaction", zap.Any("hash", txh))

	err = orm.SetBNBByMongo(hi)
	if err != nil {
		log.GetLogger().Error("mongo set BNBTransaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}

func PushH2oTransaction(from, to, contract, amount, txh, nonce, blockNumber string, code uint) error {
	txInfo, err := bsc.GetTransactionReceipt(txh, common.ScanKey)
	if err != nil {
		return err
	}
	var hi common.H20TransactionForPush
	if txInfo.Result.Status == "0x1" {
		hi.Status = "success"
	} else {
		hi.Status = "failed"
	}
	hi.To = to
	hi.From = from
	hi.Contract = contract
	hi.BlockNumber = blockNumber
	hi.Amount = amount
	hi.Nonce = nonce
	hi.Hash = txh
	hi.Currency = common.H2OCurrency

	data, err := json.Marshal(hi)
	if err != nil {
		return err
	}

	var pushIn common.PushInput
	pushIn.Data = data
	pushIn.Code = code

	pushData, err := json.Marshal(pushIn)
	if err != nil {
		return err
	}
	err = PushMongoTransaction(Producer, TopicH20ForScan, pushData)
	if err != nil {
		return err
	}
	log.GetLogger().Info("Push h2o Transaction", zap.Any("hash", txh))

	err = orm.SetH20ByMongo(hi)
	if err != nil {
		log.GetLogger().Info("mongo set h2o Transaction", zap.Error(err), zap.Any("hash", txh))
	}
	return nil
}
