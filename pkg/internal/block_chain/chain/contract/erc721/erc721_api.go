package erc721

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

var Erc721Parse abi.ABI

func init() {
	Erc721Parse, _ = abi.JSON(strings.NewReader(Erc721ABI))
}

func GetInputForCreateAsset(_propsName string, _starRating uint8, announcer common.Address) ([]byte, error) {
	input, err := Erc721Parse.Pack("createAsset", _propsName, _starRating, announcer)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func GetInputForTransfer(_to common.Address, _assetId *big.Int) ([]byte, error) {
	input, err := Erc721Parse.Pack("transfer", _to, _assetId)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func GetInputForApproval(_approved common.Address, _tokenId *big.Int) ([]byte, error) {
	input, err := Erc721Parse.Pack("approve", _approved, _tokenId)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func GetInputForTransferFrom(_from common.Address, _to common.Address, _value *big.Int) ([]byte, error) {
	input, err := Erc721Parse.Pack("transferFrom", _from, _to, _value)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func ParseTransfer(input string) (string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", err
	}

	if len(data) != 68 {
		return "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("transfer", data[4:68])
	if err != nil {
		return "", "", err
	}
	if len(res) != 2 {
		return "", "", errors.New("error:parse is wrong")
	}

	return strings.ToLower(fmt.Sprint(res[0])), fmt.Sprint(res[1]), nil
}

func ParseApproval(input string) (string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", err
	}

	if len(data) != 68 {
		return "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("approve", data[4:68])
	if err != nil {
		return "", "", err
	}
	if len(res) != 2 {
		return "", "", errors.New("error:parse is wrong")
	}

	return strings.ToLower(fmt.Sprint(res[0])), fmt.Sprint(res[1]), nil
}

func ParseTransferFrom(input string) (string, string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", "", err
	}

	if len(data) != 100 {
		return "", "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("transferFrom", data[4:100])
	if err != nil {
		return "", "", "", err
	}
	if len(res) != 3 {
		return "", "", "", errors.New("error:parse is wrong")
	}

	return strings.ToLower(fmt.Sprint(res[0])), strings.ToLower(fmt.Sprint(res[1])), fmt.Sprint(res[2]), nil
}

func MethodUnPackInputs(name string, input []byte) ([]interface{}, error) {
	method, exist := Erc721Parse.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.UnpackValues(input)
	if err != nil {
		return nil, err
	}
	return arguments, nil
}

func UnPackCreateAssetEven(data []byte) (string, string, string, string, error) {
	ev, exist := Erc721Parse.Events["CreateAsset"]
	if !exist {
		return "", "", "", "", errors.New("CreateAsset even not found")
	}
	res, err := ev.Inputs.Unpack(data)
	if err != nil {
		return "", "", "", "", err
	}

	if len(res) != 4 {
		return "", "", "", "", errors.New("result is wrong")
	}

	return strings.ToLower(fmt.Sprint(res[0])), fmt.Sprint(res[1]), fmt.Sprint(res[2]), fmt.Sprint(res[3]), nil
}
