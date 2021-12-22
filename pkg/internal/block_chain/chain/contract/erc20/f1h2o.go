package erc20

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// F1h2oABI is the input ABI used to generate the binding from.
const F1h2oABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// F1h2oBin is the compiled bytecode used for deploying new contracts.
var F1h2oBin = "0x60806040523480156200001157600080fd5b506040518060400160405280600b81526020017f463148324f20746f6b656e0000000000000000000000000000000000000000008152506040518060400160405280600581526020017f463148324f000000000000000000000000000000000000000000000000000000815250601282600390805190602001906200009892919062000232565b508160049080519060200190620000b192919062000232565b5080600560006101000a81548160ff021916908360ff160217905550505050620000f36012600a0a640ba43b7400620001ad60201b62000cf71790919060201c565b600081905550600054600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6000546040518082815260200191505060405180910390a3620002e1565b600080831415620001c257600090506200022c565b818302905081838281620001d257fe5b04146200022b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180620011746021913960400191505060405180910390fd5b5b92915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200027557805160ff1916838001178555620002a6565b82800160010185558215620002a6579182015b82811115620002a557825182559160200191906001019062000288565b5b509050620002b59190620002b9565b5090565b620002de91905b80821115620002da576000816000905550600101620002c0565b5090565b90565b610e8380620002f16000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce5671461022557806370a082311461024957806395d89b41146102a1578063a9059cbb14610324578063dd62ed3e1461038a57610093565b806306fdde0314610098578063095ea7b31461011b57806318160ddd1461018157806323b872dd1461019f575b600080fd5b6100a0610402565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e05780820151818401526020810190506100c5565b50505050905090810190601f16801561010d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101676004803603604081101561013157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506104a0565b604051808215151515815260200191505060405180910390f35b6101896104b7565b6040518082815260200191505060405180910390f35b61020b600480360360608110156101b557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506104bd565b604051808215151515815260200191505060405180910390f35b61022d61056e565b604051808260ff1660ff16815260200191505060405180910390f35b61028b6004803603602081101561025f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610581565b6040518082815260200191505060405180910390f35b6102a9610599565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102e95780820151818401526020810190506102ce565b50505050905090810190601f1680156103165780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103706004803603604081101561033a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610637565b604051808215151515815260200191505060405180910390f35b6103ec600480360360408110156103a057600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061064e565b6040518082815260200191505060405180910390f35b60038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104985780601f1061046d57610100808354040283529160200191610498565b820191906000526020600020905b81548152906001019060200180831161047b57829003601f168201915b505050505081565b60006104ad3384846106d5565b6001905092915050565b60005481565b60006104ca8484846108cc565b610563843361055e85600260008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610bf190919063ffffffff16565b6106d5565b600190509392505050565b600560009054906101000a900460ff1681565b60016020528060005260406000206000915090505481565b60048054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561062f5780601f106106045761010080835404028352916020019161062f565b820191906000526020600020905b81548152906001019060200180831161061257829003601f168201915b505050505081565b60006106443384846108cc565b6001905092915050565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561075b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180610e2b6024913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156107e1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180610d9b6022913960400191505060405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040518082815260200191505060405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610952576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180610e066025913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156109d8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180610d786023913960400191505060405180910390fd5b3073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610a5d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180610dbd6028913960400191505060405180910390fd5b610aaf81600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610bf190919063ffffffff16565b600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610b4481600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610c7490919063ffffffff16565b600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a3505050565b600082821115610c69576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525060200191505060405180910390fd5b818303905092915050565b6000818301905082811015610cf1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b92915050565b600080831415610d0a5760009050610d71565b818302905081838281610d1957fe5b0414610d70576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180610de56021913960400191505060405180910390fd5b5b9291505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220746f207468697320636f6e74726163742061646472657373536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f7745524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a265627a7a723158208e1557b70225fdffaf107d82fe115e5a26b1f242ba44a9439734e045ed17d4dc64736f6c63430005110032536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77"

// DeployF1h2o deploys a new Ethereum contract, binding an instance of F1h2o to it.
func DeployF1h2o(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *F1h2o, error) {
	parsed, err := abi.JSON(strings.NewReader(F1h2oABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(F1h2oBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &F1h2o{F1h2oCaller: F1h2oCaller{contract: contract}, F1h2oTransactor: F1h2oTransactor{contract: contract}, F1h2oFilterer: F1h2oFilterer{contract: contract}}, nil
}

// F1h2o is an auto generated Go binding around an Ethereum contract.
type F1h2o struct {
	F1h2oCaller     // Read-only binding to the contract
	F1h2oTransactor // Write-only binding to the contract
	F1h2oFilterer   // Log filterer for contract events
}

// F1h2oCaller is an auto generated read-only Go binding around an Ethereum contract.
type F1h2oCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// F1h2oTransactor is an auto generated write-only Go binding around an Ethereum contract.
type F1h2oTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// F1h2oFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type F1h2oFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// F1h2oSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type F1h2oSession struct {
	Contract     *F1h2o            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// F1h2oCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type F1h2oCallerSession struct {
	Contract *F1h2oCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// F1h2oTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type F1h2oTransactorSession struct {
	Contract     *F1h2oTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// F1h2oRaw is an auto generated low-level Go binding around an Ethereum contract.
type F1h2oRaw struct {
	Contract *F1h2o // Generic contract binding to access the raw methods on
}

// F1h2oCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type F1h2oCallerRaw struct {
	Contract *F1h2oCaller // Generic read-only contract binding to access the raw methods on
}

// F1h2oTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type F1h2oTransactorRaw struct {
	Contract *F1h2oTransactor // Generic write-only contract binding to access the raw methods on
}

// NewF1h2o creates a new instance of F1h2o, bound to a specific deployed contract.
func NewF1h2o(address common.Address, backend bind.ContractBackend) (*F1h2o, error) {
	contract, err := bindF1h2o(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &F1h2o{F1h2oCaller: F1h2oCaller{contract: contract}, F1h2oTransactor: F1h2oTransactor{contract: contract}, F1h2oFilterer: F1h2oFilterer{contract: contract}}, nil
}

// NewF1h2oCaller creates a new read-only instance of F1h2o, bound to a specific deployed contract.
func NewF1h2oCaller(address common.Address, caller bind.ContractCaller) (*F1h2oCaller, error) {
	contract, err := bindF1h2o(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &F1h2oCaller{contract: contract}, nil
}

// NewF1h2oTransactor creates a new write-only instance of F1h2o, bound to a specific deployed contract.
func NewF1h2oTransactor(address common.Address, transactor bind.ContractTransactor) (*F1h2oTransactor, error) {
	contract, err := bindF1h2o(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &F1h2oTransactor{contract: contract}, nil
}

// NewF1h2oFilterer creates a new log filterer instance of F1h2o, bound to a specific deployed contract.
func NewF1h2oFilterer(address common.Address, filterer bind.ContractFilterer) (*F1h2oFilterer, error) {
	contract, err := bindF1h2o(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &F1h2oFilterer{contract: contract}, nil
}

// bindF1h2o binds a generic wrapper to an already deployed contract.
func bindF1h2o(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(F1h2oABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_F1h2o *F1h2oRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _F1h2o.Contract.F1h2oCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_F1h2o *F1h2oRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _F1h2o.Contract.F1h2oTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_F1h2o *F1h2oRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _F1h2o.Contract.F1h2oTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_F1h2o *F1h2oCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _F1h2o.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_F1h2o *F1h2oTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _F1h2o.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_F1h2o *F1h2oTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _F1h2o.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_F1h2o *F1h2oCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_F1h2o *F1h2oSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _F1h2o.Contract.Allowance(&_F1h2o.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_F1h2o *F1h2oCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _F1h2o.Contract.Allowance(&_F1h2o.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_F1h2o *F1h2oCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_F1h2o *F1h2oSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _F1h2o.Contract.BalanceOf(&_F1h2o.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_F1h2o *F1h2oCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _F1h2o.Contract.BalanceOf(&_F1h2o.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_F1h2o *F1h2oCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_F1h2o *F1h2oSession) Decimals() (uint8, error) {
	return _F1h2o.Contract.Decimals(&_F1h2o.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_F1h2o *F1h2oCallerSession) Decimals() (uint8, error) {
	return _F1h2o.Contract.Decimals(&_F1h2o.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_F1h2o *F1h2oCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_F1h2o *F1h2oSession) Name() (string, error) {
	return _F1h2o.Contract.Name(&_F1h2o.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_F1h2o *F1h2oCallerSession) Name() (string, error) {
	return _F1h2o.Contract.Name(&_F1h2o.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_F1h2o *F1h2oCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_F1h2o *F1h2oSession) Symbol() (string, error) {
	return _F1h2o.Contract.Symbol(&_F1h2o.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_F1h2o *F1h2oCallerSession) Symbol() (string, error) {
	return _F1h2o.Contract.Symbol(&_F1h2o.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_F1h2o *F1h2oCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _F1h2o.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_F1h2o *F1h2oSession) TotalSupply() (*big.Int, error) {
	return _F1h2o.Contract.TotalSupply(&_F1h2o.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_F1h2o *F1h2oCallerSession) TotalSupply() (*big.Int, error) {
	return _F1h2o.Contract.TotalSupply(&_F1h2o.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_F1h2o *F1h2oTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.contract.Transact(opts, "approve", _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_F1h2o *F1h2oSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.Approve(&_F1h2o.TransactOpts, _spender, _value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _value) returns(bool)
func (_F1h2o *F1h2oTransactorSession) Approve(_spender common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.Approve(&_F1h2o.TransactOpts, _spender, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.Transfer(&_F1h2o.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.Transfer(&_F1h2o.TransactOpts, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.contract.Transact(opts, "transferFrom", _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.TransferFrom(&_F1h2o.TransactOpts, _from, _to, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _value) returns(bool _success)
func (_F1h2o *F1h2oTransactorSession) TransferFrom(_from common.Address, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _F1h2o.Contract.TransferFrom(&_F1h2o.TransactOpts, _from, _to, _value)
}

// F1h2oApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the F1h2o contract.
type F1h2oApprovalIterator struct {
	Event *F1h2oApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *F1h2oApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(F1h2oApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(F1h2oApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *F1h2oApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *F1h2oApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// F1h2oApproval represents a Approval event raised by the F1h2o contract.
type F1h2oApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_F1h2o *F1h2oFilterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _spender []common.Address) (*F1h2oApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _F1h2o.contract.FilterLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return &F1h2oApprovalIterator{contract: _F1h2o.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_F1h2o *F1h2oFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *F1h2oApproval, _owner []common.Address, _spender []common.Address) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _spenderRule []interface{}
	for _, _spenderItem := range _spender {
		_spenderRule = append(_spenderRule, _spenderItem)
	}

	logs, sub, err := _F1h2o.contract.WatchLogs(opts, "Approval", _ownerRule, _spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(F1h2oApproval)
				if err := _F1h2o.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _spender, uint256 _value)
func (_F1h2o *F1h2oFilterer) ParseApproval(log types.Log) (*F1h2oApproval, error) {
	event := new(F1h2oApproval)
	if err := _F1h2o.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// F1h2oTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the F1h2o contract.
type F1h2oTransferIterator struct {
	Event *F1h2oTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *F1h2oTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(F1h2oTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(F1h2oTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *F1h2oTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *F1h2oTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// F1h2oTransfer represents a Transfer event raised by the F1h2o contract.
type F1h2oTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_F1h2o *F1h2oFilterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address) (*F1h2oTransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _F1h2o.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return &F1h2oTransferIterator{contract: _F1h2o.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_F1h2o *F1h2oFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *F1h2oTransfer, _from []common.Address, _to []common.Address) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}

	logs, sub, err := _F1h2o.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(F1h2oTransfer)
				if err := _F1h2o.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 _value)
func (_F1h2o *F1h2oFilterer) ParseTransfer(log types.Log) (*F1h2oTransfer, error) {
	event := new(F1h2oTransfer)
	if err := _F1h2o.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
