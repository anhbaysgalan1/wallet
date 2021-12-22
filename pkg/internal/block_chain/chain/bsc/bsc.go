package bsc

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strconv"
)

func GetBlockInfo(number uint64, key string) (BlockInfoBsc, error) {
	res, err := requestExec(EthBlock, key, "&tag="+hexutil.EncodeUint64(number)+"&boolean=true")
	if err != nil {
		return BlockInfoBsc{}, err
	}
	var bi BlockInfoBsc
	err = json.Unmarshal(res, &bi)
	if err != nil {
		return BlockInfoBsc{}, err
	}
	return bi, nil
}

func GetTransactionReceipt(txh, key string) (TransactionReceiptBsc, error) {
	res, err := requestExec(EthTransactionReceipt, key, "&txhash="+txh+"&tag=latest")
	if err != nil {
		return TransactionReceiptBsc{}, err
	}

	var tr TransactionReceiptBsc
	err = json.Unmarshal(res, &tr)
	if err != nil {
		return TransactionReceiptBsc{}, err
	}
	return tr, nil
}

func GetGasTracker(key string) (map[string]uint64, error) {
	res, err := requestExec(EthGasTracker, key, "")
	if err != nil {
		return nil, err
	}

	var gt GasTracker
	err = json.Unmarshal(res, &gt)
	if err != nil {
		return nil, err
	}

	m := make(map[string]uint64)
	fastGasPrice, err := strconv.ParseUint(gt.Result.SafeGasPrice, 0, 0)
	if err != nil {
		return nil, err
	}

	if fastGasPrice > 1000 {
		return nil, errors.New("gasPrice too big")
	}

	m["FastGasPrice"] = fastGasPrice
	return m, nil

	/*
		lastBlock, err := strconv.ParseUint(gt.Result.LastBlock, 0, 0)
		if err != nil {
			return nil, err
		}

		safeGasPrice, err := strconv.ParseUint(gt.Result.SafeGasPrice, 0, 0)
		if err != nil {
			return nil, err
		}

		proposeGasPrice, err := strconv.ParseUint(gt.Result.ProposeGasPrice, 0, 0)
		if err != nil {
			return nil, err
		}
		m["LastBlock"] = lastBlock
		m["SafeGasPrice"] = safeGasPrice
		m["ProposeGasPrice"] = proposeGasPrice
	*/
}

func GetGasPrice(key string) (uint64, error) {
	res, err := requestExec(EthGasPrice, key, "")
	if err != nil {
		return 0, err
	}

	var tcr JsonResult
	err = json.Unmarshal(res, &tcr)
	if err != nil {
		return 0, err
	}

	gp, err := hexutil.DecodeUint64(tcr.Result)
	if err != nil {
		return 0, err
	}
	return gp, nil
}

func GetAccountNonce(addr, key string) (uint64, error) {
	res, err := requestExec(EthTransactionCount, key, "&address="+addr+"&tag=latest")
	if err != nil {
		return 0, err
	}

	var tcr JsonResult
	err = json.Unmarshal(res, &tcr)
	if err != nil {
		return 0, err
	}

	noc, err := hexutil.DecodeUint64(tcr.Result)
	if err != nil {
		return 0, err
	}
	return noc, nil
}

func SendRawTransaction(raw, key string) (string, error) {
	res, err := requestExec(EthSendRawTransaction, key, "&hex="+raw)
	if err != nil {
		return "", err
	}

	var tcr JsonResult
	err = json.Unmarshal(res, &tcr)
	if err != nil {
		var errResult ErrorForEthereumApi
		err1 := json.Unmarshal(res, &errResult)
		if err1 != nil {
			return "", err1
		}
		//  {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"nonce too low"}}
		return "", errors.New(errResult.ErrorResult.Message)
	}
	return tcr.Result, nil
}
