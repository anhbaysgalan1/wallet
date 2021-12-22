package transfer

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"tp_wallet/internal/block_chain/chain/contract/erc1155"
	"tp_wallet/internal/block_chain/chain/contract/erc20"
	"tp_wallet/internal/block_chain/chain/contract/erc721"
)

//WithdrawalRacingBoatForBsc 提现赛艇NFT并直接下发给用户第三方钱包地址,返回交易hash，提币from地址nonce值，错误信息
func WithdrawalRacingBoatForBsc(data InputForTransfer, netType string, NftID *big.Int) (string, uint64, error) {
	//
	//if data.GasPrice > 50 {
	//	return "", 0, errors.New("gasPrice to high")
	//}
	//
	//if data.GasLimit < 200000 {
	//	return "", 0, errors.New("GasLimit to low")
	//}

	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 {
		return "", 0, errors.New("params is wrong")
	}

	//用户第三方接收钱包地址、赛艇nftID
	input, err := erc721.GetInputForTransfer(common.HexToAddress(data.ToAddress), NftID)
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1500000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//WithdrawalMaterialForBsc 提现材料NFT并直接下发给指定钱包地址,返回交易hash，提币from地址nonce值，错误信息
func WithdrawalMaterialForBsc(data InputForTransfer, netType string, NftID *big.Int, _amount *big.Int) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 {
		return "", 0, errors.New("params is wrong")
	}

	//用户from地址、to地址、材料nftID、材料数量
	input, err := erc1155.GetInputForTransfer(
		common.HexToAddress(data.FromAddress),
		common.HexToAddress(data.ToAddress),
		NftID,
		_amount)
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1000000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//CreateAssetMaterialForBsc 生产材料NFT并直接下发给用户第三方钱包地址,返回交易hash，提币from地址nonce值，错误信息
func CreateAssetMaterialForBsc(data InputForTransfer, netType string, _amount *big.Int, _name string) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 ||
		len(_name) == 0 {
		return "", 0, errors.New("params is wrong")
	}

	//材料数量、接收钱包地址、材料名称
	input, err := erc1155.GetInputForCreateAsset(_amount, common.HexToAddress(data.ToAddress), _name)
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1000000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//ExpandMaterialForBsc 提现材料NFT并直接下发给指定钱包地址,返回交易hash，提币from地址nonce值，错误信息
func ExpandMaterialForBsc(data InputForTransfer, netType string, NftID *big.Int, _amount *big.Int) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 {
		return "", 0, errors.New("params is wrong")
	}

	//材料nftID、材料数量
	input, err := erc1155.GetInputForExpand(
		NftID,
		_amount)
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1000000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//CreateAssetRacingBoatForBsc 生产赛艇NFT并直接下发给用户第三方钱包地址,返回交易hash，提币from地址nonce值，错误信息
func CreateAssetRacingBoatForBsc(data InputForTransfer, netType string, _propsName string, _starRating uint8) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 ||
		len(_propsName) == 0 {
		return "", 0, errors.New("params is wrong")
	}

	//赛艇名称、星级、接收钱包地址
	input, err := erc721.GetInputForCreateAsset(_propsName, _starRating, common.HexToAddress(data.ToAddress))
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1000000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//Withdrawal20TokenForBsc 提取合约代币,返回交易hash，提币from地址nonce值，错误信息
func Withdrawal20TokenForBsc(data InputForTransfer, netType string) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 ||
		len(data.ContractAddress) != 42 {
		return "", 0, errors.New("params is wrong")
	}

	input, err := erc20.GetInputForTransfer(common.HexToAddress(data.ToAddress), data.Amount)
	if err != nil {
		return "", 0, err
	}

	var cli *ethclient.Client
	var gasP uint64
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		1000000,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

//WithdrawalBNBForBsc 提取赛艇原生代币,返回交易hash，提币from地址nonce值，错误信息
func WithdrawalBNBForBsc(data InputForBNBTransfer, netType string) (string, uint64, error) {
	if len(data.FromAddress) != 42 ||
		len(data.Private) == 0 ||
		len(data.ToAddress) != 42 {
		return "", 0, errors.New("params is wrong")
	}

	var cli *ethclient.Client
	var gasP uint64
	var err error
	switch netType {
	case ChainNetWorkTest:
		cli, err = ethclient.Dial(TestNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 10
	case ChainNetWorkRelease:
		cli, err = ethclient.Dial(MainNetRpcUrl1)
		if err != nil {
			return "", 0, err
		}
		gasP = 5
	default:
		if err != nil {
			return "", 0, errors.New("netType is wrong: " + netType)
		}
	}

	hash, err := SendRawTransactionByRpc(
		data.Private,
		cli,
		common.HexToAddress(data.ToAddress),
		data.Amount,
		21000,
		data.Nonce,
		nil,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(gasP)))
	if err != nil {
		return "", 0, err
	}
	return hash, data.Nonce, nil
}

/*
//WithdrawalBNBForBsc2 提取赛艇原生代币,返回交易hash，提币总账户nonce值，错误信息
func WithdrawalBNBForBsc2(data InputForBNBTransfer, key string) (string, uint64, error) {
	ops, err := NewTransferOpt(data.Private, 56)
	if err != nil {
		return "", 0, err
	}

	if data.GasPrice > 50 {
		return "", 0, errors.New("gasPrice to high")
	}

	if data.GasLimit < 200000 {
		return "", 0, errors.New("GasLimit to low")
	}

	hash, raw, err := RawTransactionToString(
		ops,
		common.HexToAddress(data.ToAddress),
		data.Amount,
		data.GasLimit,
		data.Nonce,
		nil,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(data.GasPrice)))
	if err != nil {
		return "", 0, err
	}

	res, err := bsc.SendRawTransaction(raw, key)
	if err != nil {
		return "", 0, err
	}

	if strings.ToLower(res) != strings.ToLower(hash.String()) {
		return "", 0, errors.New("error: send Bsc transaction for Bsc mainnet is wrong!--res:" + res)
	}
	//重复发送一遍交易，预防之前的漏发问题
	bsc.SendRawTransaction(raw, key)
	return res, data.Nonce, nil
}

//Withdrawal20TokenForBsc2 提取合约代币,返回交易hash，提币总账户nonce值，错误信息
func Withdrawal20TokenForBsc2(data InputForTransfer, key string) (string, uint64, error) {
	ops, err := NewTransferOpt(data.Private, 56)
	if err != nil {
		return "", 0, err
	}

	if data.GasPrice > 50 {
		return "", 0, errors.New("gasPrice to high")
	}

	if data.GasLimit < 200000 {
		return "", 0, errors.New("GasLimit to low")
	}

	input, err := erc20.GetInputForTransfer(common.HexToAddress(data.ToAddress), data.Amount)
	if err != nil {
		return "", 0, err
	}

	hash, raw, err := RawTransactionToString(
		ops,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		data.GasLimit,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(data.GasPrice)))
	if err != nil {
		return "", 0, err
	}

	res, err := bsc.SendRawTransaction(raw, key)
	if err != nil {
		return "", 0, err
	}

	if strings.ToLower(res) != strings.ToLower(hash.String()) {
		return "", 0, errors.New("error: send Bsc transaction for Bsc mainnet is wrong!--res:" + res)
	}
	//重复发送一遍交易，预防之前的漏发问题
	bsc.SendRawTransaction(raw, key)
	return res, data.Nonce, nil
}

//CreateAssetRacingBoatForBsc2 生产赛艇NFT并直接下发给用户第三方钱包地址,返回交易hash，提币总账户nonce值，错误信息
func CreateAssetRacingBoatForBsc2(data InputForTransfer, key string, _propsName string, _starRating uint8) (string, uint64, error) {
	ops, err := NewTransferOpt(data.Private, 56)
	if err != nil {
		return "", 0, err
	}

	if data.GasPrice > 50 {
		return "", 0, errors.New("gasPrice to high")
	}

	if data.GasLimit < 200000 {
		return "", 0, errors.New("GasLimit to low")
	}

	//赛艇名称、星级、接收钱包地址
	input, err := erc721.GetInputForCreateAsset(_propsName, _starRating, common.HexToAddress(data.ToAddress))
	if err != nil {
		return "", 0, err
	}

	hash, raw, err := RawTransactionToString(
		ops,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		data.GasLimit,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(data.GasPrice)))
	if err != nil {
		return "", 0, err
	}

	res, err := bsc.SendRawTransaction(raw, key)
	if err != nil {
		return "", 0, err
	}

	if strings.ToLower(res) != strings.ToLower(hash.String()) {
		return "", 0, errors.New("error: send bsc transaction for Bsc mainnet is wrong!--res:" + res)
	}
	//重复发送一遍交易，预防之前的漏发问题
	bsc.SendRawTransaction(raw, key)
	return res, data.Nonce, nil
}

//WithdrawalRacingBoatForBsc2 生产赛艇NFT并直接下发给用户第三方钱包地址,返回交易hash，提币总账户nonce值，错误信息
func WithdrawalRacingBoatForBsc2(data InputForTransfer, key string, NftID *big.Int) (string, uint64, error) {
	ops, err := NewTransferOpt(data.Private, 56)
	if err != nil {
		return "", 0, err
	}

	if data.GasPrice > 50 {
		return "", 0, errors.New("gasPrice to high")
	}

	if data.GasLimit < 200000 {
		return "", 0, errors.New("GasLimit to low")
	}

	//用户第三方接收钱包地址、赛艇nftID
	input, err := erc721.GetInputForTransfer(common.HexToAddress(data.ToAddress), NftID)
	if err != nil {
		return "", 0, err
	}

	hash, raw, err := RawTransactionToString(
		ops,
		common.HexToAddress(data.ContractAddress),
		new(big.Int).SetInt64(0),
		data.GasLimit,
		data.Nonce,
		input,
		new(big.Int).Mul(big.NewInt(1e+9), new(big.Int).SetUint64(data.GasPrice)))
	if err != nil {
		return "", 0, err
	}

	res, err := bsc.SendRawTransaction(raw, key)
	if err != nil {
		return "", 0, err
	}

	if strings.ToLower(res) != strings.ToLower(hash.String()) {
		return "", 0, errors.New("error: send bsc transaction for Bsc mainnet is wrong!--res:" + res)
	}
	//重复发送一遍交易，预防之前的漏发问题
	bsc.SendRawTransaction(raw, key)
	return res, data.Nonce, nil
}
*/
