package erc1155

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

var Erc1155Parse abi.ABI

func init() {
	Erc1155Parse, _ = abi.JSON(strings.NewReader(MateralABI))
}

func GetInputForCreateAsset(amount *big.Int, announcer common.Address, materialName string) ([]byte, error) {
	input, err := Erc1155Parse.Pack("createAsset", amount, announcer, materialName)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func ParseCreateAsset(input string) (string, string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", "", err
	}

	if len(data) != 164 {
		return "", "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("createAsset", data[4:164])
	if err != nil {
		return "", "", "", err
	}
	if len(res) != 3 {
		return "", "", "", errors.New("error:parse is wrong")
	}

	return fmt.Sprint(res[0]), strings.ToLower(fmt.Sprint(res[1])), fmt.Sprint(res[2]), nil
}

func GetInputForTransfer(from common.Address, to common.Address, id *big.Int, amount *big.Int) ([]byte, error) {
	input, err := Erc1155Parse.Pack("safeTransferFrom", from, to, id, amount)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func ParseTransfer(input string) (string, string, string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", "", "", err
	}

	if len(data) != 132 {
		return "", "", "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("safeTransferFrom", data[4:132])
	if err != nil {
		return "", "", "", "", err
	}
	if len(res) != 4 {
		return "", "", "", "", errors.New("error:parse is wrong")
	}

	return strings.ToLower(fmt.Sprint(res[0])), strings.ToLower(fmt.Sprint(res[1])), fmt.Sprint(res[2]), fmt.Sprint(res[3]), nil
}

func GetInputForExpand(id *big.Int, amount *big.Int) ([]byte, error) {
	input, err := Erc1155Parse.Pack("expand", id, amount)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func ParseExpand(input string) (string, string, error) {
	data, err := hexutil.Decode(input)
	if err != nil {
		return "", "", err
	}

	if len(data) != 68 {
		return "", "", errors.New("error:parse is wrong")
	}

	res, err := MethodUnPackInputs("expand", data[4:68])
	if err != nil {
		return "", "", err
	}
	if len(res) != 2 {
		return "", "", errors.New("error:parse is wrong")
	}

	return fmt.Sprint(res[0]), fmt.Sprint(res[1]), nil
}

func MethodUnPackInputs(name string, input []byte) ([]interface{}, error) {
	method, exist := Erc1155Parse.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.UnpackValues(input)
	if err != nil {
		return nil, err
	}
	return arguments, nil
}

// UnPackCreateAssetEven 返回：address announcer,uint256 id,string materialName,uint256 amount
func UnPackCreateAssetEven(data []byte) (string, string, string, string, error) {
	ev, exist := Erc1155Parse.Events["CreateAsset"]
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
