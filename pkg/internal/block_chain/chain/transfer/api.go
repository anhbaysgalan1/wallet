package transfer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

// RawTransactionToString Normal transaction
func RawTransactionToString(opts *bind.TransactOpts,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	nonce uint64,
	data []byte,
	gasPrice *big.Int) (common.Hash, string, error) {

	// Create the transaction, sign it and schedule it for execution
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	if opts.Signer == nil {
		return common.Hash{}, "", errors.New("no signer to authorize the transaction with")
	}

	//sign transaction
	signedTx, err := opts.Signer(opts.From, rawTx)
	if err != nil {
		return common.Hash{}, "", err
	}

	txB, err := signedTx.MarshalBinary()
	if err != nil {
		return common.Hash{}, "", err
	}

	return signedTx.Hash(), hexutil.Encode(txB), nil
}

// SendRawTransactionByRpc Normal transaction
func SendRawTransactionByRpc(
	private string,
	cli *ethclient.Client,
	to common.Address,
	amount *big.Int,
	gasLimit uint64,
	nonce uint64,
	data []byte,
	gasPrice *big.Int) (string, error) {

	opts, err := NewTransferOpt(private, 97)
	if err != nil {
		return "", err
	}

	// Create the transaction, sign it and schedule it for execution
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	if opts.Signer == nil {
		return "", errors.New("no signer to authorize the transaction with")
	}

	//sign transaction
	signedTx, err := opts.Signer(opts.From, rawTx)
	if err != nil {
		return "", err
	}

	//send transaction
	err = cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	//time.Sleep(time.Second / 5)
	//cli.SendTransaction(context.Background(), signedTx)
	return signedTx.Hash().String(), nil
}

func NewTransferOpt(key string, Code uint64) (*bind.TransactOpts, error) {
	pri, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	opt, err := bind.NewKeyedTransactorWithChainID(pri, new(big.Int).SetUint64(Code))
	if err != nil {
		return nil, err
	}
	return opt, nil
}

func GetAccountNonce(addr, netType string) (uint64, error) {
	var url string
	switch netType {
	case ChainNetWorkTest:
		url = UrlTestNet
	case ChainNetWorkRelease:
		url = UrlMainNet
	default:
		return 0, errors.New("netType is wrong: " + netType)
	}
	res, err := requestExec(EthTransactionCount, url, "&address="+addr+"&tag=latest")
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

func requestExec(method, net, param string) ([]byte, error) {
	time.Sleep(time.Second / 4)
	url := net + method + param + "&apikey=" + BscApi
	tl := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tl}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
