package erc20

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strings"
)

var Erc20Parse abi.ABI

func init() {
	Erc20Parse, _ = abi.JSON(strings.NewReader(F1h2oABI))
}

func MethodUnPackInputs(name string, input []byte) ([]interface{}, error) {
	method, exist := Erc20Parse.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.UnpackValues(input)
	if err != nil {
		return nil, err
	}
	return arguments, nil
}

func MethodPackInputs(name string, to common.Address, value *big.Int) ([]byte, []byte, error) {
	method, exist := Erc20Parse.Methods[name]
	if !exist {
		return nil, nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.Pack(to, value)
	if err != nil {
		return nil, nil, err
	}
	return method.ID, arguments, nil
}

func GetInputForTransfer(to common.Address, value *big.Int) ([]byte, error) {
	input, err := Erc20Parse.Pack("transfer", to, value)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func GetInputForApproval(_approved common.Address, _value *big.Int) ([]byte, error) {
	input, err := Erc20Parse.Pack("approve", _approved, _value)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func GetInputForTransferFrom(_from common.Address, _to common.Address, _value *big.Int) ([]byte, error) {
	input, err := Erc20Parse.Pack("transferFrom", _from, _to, _value)
	if err != nil {
		return nil, err
	}
	return input, nil
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
