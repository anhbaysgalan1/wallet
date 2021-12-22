package erc721

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

// Erc721ABI is the input ABI used to generate the binding from.
const Erc721ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_propsName\",\"type\":\"string\"},{\"name\":\"_starRating\",\"type\":\"uint8\"},{\"name\":\"announcer\",\"type\":\"address\"}],\"name\":\"createAsset\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"_approved\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_approved\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_assetId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gameName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"_balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_assetId\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"getAssetsByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAssetNumber\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_assetId\",\"type\":\"uint256\"}],\"name\":\"getAsset\",\"outputs\":[{\"name\":\"_gameName\",\"type\":\"string\"},{\"name\":\"_propsName\",\"type\":\"string\"},{\"name\":\"_starRating\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"announcer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_propsName\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"_starRating\",\"type\":\"uint8\"}],\"name\":\"CreateAsset\",\"type\":\"event\"}]"

// Erc721Bin is the compiled bytecode used for deploying new contracts.
var Erc721Bin = "0x60806040526040805190810160405280600c81526020017f46312048324f2053504545440000000000000000000000000000000000000000815250600190805190602001906200005192919062000098565b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555062000147565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000db57805160ff19168380011785556200010c565b828001600101855582156200010c579182015b828111156200010b578251825591602001919060010190620000ee565b5b5090506200011b91906200011f565b5090565b6200014491905b808211156200014057600081600090555060010162000126565b5090565b90565b61162780620001576000396000f3006080604052600436106100af576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630324143f146100b4578063081812fc1461015e578063095ea7b3146101cb57806323b872dd1461020b578063473bc2231461026b5780636352211e146102fb57806370a0823114610368578063a9059cbb146103bf578063d5b3b67d146103ff578063df15f11114610497578063eac8f5b8146104c2575b600080fd5b3480156100c057600080fd5b50610148600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803560ff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105e1565b6040518082815260200191505060405180910390f35b34801561016a57600080fd5b506101896004803603810190808035906020019092919050505061094e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610209600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610a87565b005b610269600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610c0e565b005b34801561027757600080fd5b50610280610df0565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102c05780820151818401526020810190506102a5565b50505050905090810190601f1680156102ed5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561030757600080fd5b5061032660048036038101908080359060200190929190505050610e8e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561037457600080fd5b506103a9600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ecb565b6040518082815260200191505060405180910390f35b6103fd600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610f14565b005b34801561040b57600080fd5b50610440600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ff9565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610483578082015181840152602081019050610468565b505050509050019250505060405180910390f35b3480156104a357600080fd5b506104ac61112b565b6040518082815260200191505060405180910390f35b3480156104ce57600080fd5b506104ed60048036038101908080359060200190929190505050611138565b6040518080602001806020018460ff1660ff168152602001838103835286818151815260200191508051906020019080838360005b8381101561053d578082015181840152602081019050610522565b50505050905090810190601f16801561056a5780820380516001836020036101000a031916815260200191505b50838103825285818151815260200191508051906020019080838360005b838110156105a3578082015181840152602081019050610588565b50505050905090810190601f1680156105d05780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b60006105eb611539565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156106b1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f41737365743a6d73672e73656e646572206e6f74206973206f776e657221000081525060200191505060405180910390fd5b60058560ff161115151561072d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f41737365743a5374617220726174696e6720697320746f6f206869677421000081525060200191505060405180910390fd5b85826000018190525084826020019060ff16908160ff16815250506001600283908060018154018082558091505090600182039060005260206000209060020201600090919290919091506000820151816000019080519060200190610794929190611556565b5060208201518160010160006101000a81548160ff021916908360ff1602179055505050039050836003600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600081548092919060010191905055507f50484cca76d2106d636492405f185725c2cafdbe8855cb341deb4b882898a79084828888604051808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001848152602001806020018360ff1660ff168152602001828103825284818151815260200191508051906020019080838360005b838110156109055780820151818401526020810190506108ea565b50505050905090810190601f1680156109325780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a180925050509392505050565b6000806002805490501115156109cc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f41737365743a4173736574206973206e756c6c2100000000000000000000000081525060200191505060405180910390fd5b6001600280549050038211151515610a4c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f41737365743a6173736574496420697320746f6f20686967742100000000000081525060200191505060405180910390fd5b6005600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6003600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610b5d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4552433732313a6d73672e73656e646572206e6f74206973206f776e6572210081525060200191505060405180910390fd5b816005600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550808273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b6003600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141515610ce4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f4552433732313a66726f6d206e6f74206973206f776e6572210000000000000081525060200191505060405180910390fd5b6005600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610de0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f4552433732313a6d73672e73656e646572206e6f7420697320617070726f766581526020017f642100000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b610deb8383836113b2565b505050565b60018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610e865780601f10610e5b57610100808354040283529160200191610e86565b820191906000526020600020905b815481529060010190602001808311610e6957829003601f168201915b505050505081565b60006003600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050919050565b6000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b6003600082815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610fea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f4552433732313a6d73672e73656e646572206e6f74206973206f776e6572210081525060200191505060405180910390fd5b610ff53383836113b2565b5050565b606080600080600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205460405190808252806020026020018201604052801561106d5781602001602082028038833980820191505090505b50925060009150600090505b600280549050811015611120578473ffffffffffffffffffffffffffffffffffffffff166003600083815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415611113578083838151811015156110fc57fe5b906020019060200201818152505081806001019250505b8080600101915050611079565b829350505050919050565b6000600280549050905090565b60608060008060006002805490501115156111bb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f41737365743a4173736574206973206e756c6c2100000000000000000000000081525060200191505060405180910390fd5b600160028054905003851115151561123b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f41737365743a6173736574496420697320746f6f20686967742100000000000081525060200191505060405180910390fd5b60028581548110151561124a57fe5b9060005260206000209060020201905060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156112f05780601f106112c5576101008083540402835291602001916112f0565b820191906000526020600020905b8154815290600101906020018083116112d357829003601f168201915b50505050509350806000018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561138f5780601f106113645761010080835404028352916020019161138f565b820191906000526020600020905b81548152906001019060200180831161137257829003601f168201915b505050505092508060010160009054906101000a900460ff169150509193909250565b600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008154809291906001019190505550600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000815480929190600190039190505550816003600083815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506005600082815260200190815260200160002060006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055808273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef60405160405180910390a4505050565b604080519081016040528060608152602001600060ff1681525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061159757805160ff19168380011785556115c5565b828001600101855582156115c5579182015b828111156115c45782518255916020019190600101906115a9565b5b5090506115d291906115d6565b5090565b6115f891905b808211156115f45760008160009055506001016115dc565b5090565b905600a165627a7a72305820fb1fb966a415750e9372c7c9c95eb50872fd83ec0c5ec90778547b4ba8a8ac2c0029"

// DeployErc721 deploys a new Ethereum contract, binding an instance of Erc721 to it.
func DeployErc721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Erc721, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc721ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Erc721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// Erc721 is an auto generated Go binding around an Ethereum contract.
type Erc721 struct {
	Erc721Caller     // Read-only binding to the contract
	Erc721Transactor // Write-only binding to the contract
	Erc721Filterer   // Log filterer for contract events
}

// Erc721Caller is an auto generated read-only Go binding around an Ethereum contract.
type Erc721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc721Session struct {
	Contract     *Erc721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc721CallerSession struct {
	Contract *Erc721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Erc721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc721TransactorSession struct {
	Contract     *Erc721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721Raw is an auto generated low-level Go binding around an Ethereum contract.
type Erc721Raw struct {
	Contract *Erc721 // Generic contract binding to access the raw methods on
}

// Erc721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc721CallerRaw struct {
	Contract *Erc721Caller // Generic read-only contract binding to access the raw methods on
}

// Erc721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc721TransactorRaw struct {
	Contract *Erc721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewErc721 creates a new instance of Erc721, bound to a specific deployed contract.
func NewErc721(address common.Address, backend bind.ContractBackend) (*Erc721, error) {
	contract, err := bindErc721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// NewErc721Caller creates a new read-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Caller(address common.Address, caller bind.ContractCaller) (*Erc721Caller, error) {
	contract, err := bindErc721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Caller{contract: contract}, nil
}

// NewErc721Transactor creates a new write-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Transactor(address common.Address, transactor bind.ContractTransactor) (*Erc721Transactor, error) {
	contract, err := bindErc721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Transactor{contract: contract}, nil
}

// NewErc721Filterer creates a new log filterer instance of Erc721, bound to a specific deployed contract.
func NewErc721Filterer(address common.Address, filterer bind.ContractFilterer) (*Erc721Filterer, error) {
	contract, err := bindErc721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc721Filterer{contract: contract}, nil
}

// bindErc721 binds a generic wrapper to an already deployed contract.
func bindErc721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc721ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.Erc721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 _balance)
func (_Erc721 *Erc721Caller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 _balance)
func (_Erc721 *Erc721Session) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 _balance)
func (_Erc721 *Erc721CallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, _owner)
}

// GameName is a free data retrieval call binding the contract method 0x473bc223.
//
// Solidity: function gameName() view returns(string)
func (_Erc721 *Erc721Caller) GameName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "gameName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GameName is a free data retrieval call binding the contract method 0x473bc223.
//
// Solidity: function gameName() view returns(string)
func (_Erc721 *Erc721Session) GameName() (string, error) {
	return _Erc721.Contract.GameName(&_Erc721.CallOpts)
}

// GameName is a free data retrieval call binding the contract method 0x473bc223.
//
// Solidity: function gameName() view returns(string)
func (_Erc721 *Erc721CallerSession) GameName() (string, error) {
	return _Erc721.Contract.GameName(&_Erc721.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 _tokenId) view returns(address _approved)
func (_Erc721 *Erc721Caller) GetApproved(opts *bind.CallOpts, _tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getApproved", _tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 _tokenId) view returns(address _approved)
func (_Erc721 *Erc721Session) GetApproved(_tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, _tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 _tokenId) view returns(address _approved)
func (_Erc721 *Erc721CallerSession) GetApproved(_tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, _tokenId)
}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 _assetId) view returns(string _gameName, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721Caller) GetAsset(opts *bind.CallOpts, _assetId *big.Int) (struct {
	GameName   string
	PropsName  string
	StarRating uint8
}, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getAsset", _assetId)

	outstruct := new(struct {
		GameName   string
		PropsName  string
		StarRating uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameName = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.PropsName = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.StarRating = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 _assetId) view returns(string _gameName, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721Session) GetAsset(_assetId *big.Int) (struct {
	GameName   string
	PropsName  string
	StarRating uint8
}, error) {
	return _Erc721.Contract.GetAsset(&_Erc721.CallOpts, _assetId)
}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 _assetId) view returns(string _gameName, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721CallerSession) GetAsset(_assetId *big.Int) (struct {
	GameName   string
	PropsName  string
	StarRating uint8
}, error) {
	return _Erc721.Contract.GetAsset(&_Erc721.CallOpts, _assetId)
}

// GetAssetNumber is a free data retrieval call binding the contract method 0xdf15f111.
//
// Solidity: function getAssetNumber() view returns(uint256)
func (_Erc721 *Erc721Caller) GetAssetNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getAssetNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetNumber is a free data retrieval call binding the contract method 0xdf15f111.
//
// Solidity: function getAssetNumber() view returns(uint256)
func (_Erc721 *Erc721Session) GetAssetNumber() (*big.Int, error) {
	return _Erc721.Contract.GetAssetNumber(&_Erc721.CallOpts)
}

// GetAssetNumber is a free data retrieval call binding the contract method 0xdf15f111.
//
// Solidity: function getAssetNumber() view returns(uint256)
func (_Erc721 *Erc721CallerSession) GetAssetNumber() (*big.Int, error) {
	return _Erc721.Contract.GetAssetNumber(&_Erc721.CallOpts)
}

// GetAssetsByAddress is a free data retrieval call binding the contract method 0xd5b3b67d.
//
// Solidity: function getAssetsByAddress(address _account) view returns(uint256[])
func (_Erc721 *Erc721Caller) GetAssetsByAddress(opts *bind.CallOpts, _account common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getAssetsByAddress", _account)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAssetsByAddress is a free data retrieval call binding the contract method 0xd5b3b67d.
//
// Solidity: function getAssetsByAddress(address _account) view returns(uint256[])
func (_Erc721 *Erc721Session) GetAssetsByAddress(_account common.Address) ([]*big.Int, error) {
	return _Erc721.Contract.GetAssetsByAddress(&_Erc721.CallOpts, _account)
}

// GetAssetsByAddress is a free data retrieval call binding the contract method 0xd5b3b67d.
//
// Solidity: function getAssetsByAddress(address _account) view returns(uint256[])
func (_Erc721 *Erc721CallerSession) GetAssetsByAddress(_account common.Address) ([]*big.Int, error) {
	return _Erc721.Contract.GetAssetsByAddress(&_Erc721.CallOpts, _account)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address _owner)
func (_Erc721 *Erc721Caller) OwnerOf(opts *bind.CallOpts, _tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "ownerOf", _tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address _owner)
func (_Erc721 *Erc721Session) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, _tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address _owner)
func (_Erc721 *Erc721CallerSession) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _approved, uint256 _tokenId) payable returns()
func (_Erc721 *Erc721Transactor) Approve(opts *bind.TransactOpts, _approved common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "approve", _approved, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _approved, uint256 _tokenId) payable returns()
func (_Erc721 *Erc721Session) Approve(_approved common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, _approved, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _approved, uint256 _tokenId) payable returns()
func (_Erc721 *Erc721TransactorSession) Approve(_approved common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, _approved, _tokenId)
}

// CreateAsset is a paid mutator transaction binding the contract method 0x0324143f.
//
// Solidity: function createAsset(string _propsName, uint8 _starRating, address announcer) returns(uint256)
func (_Erc721 *Erc721Transactor) CreateAsset(opts *bind.TransactOpts, _propsName string, _starRating uint8, announcer common.Address) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "createAsset", _propsName, _starRating, announcer)
}

// CreateAsset is a paid mutator transaction binding the contract method 0x0324143f.
//
// Solidity: function createAsset(string _propsName, uint8 _starRating, address announcer) returns(uint256)
func (_Erc721 *Erc721Session) CreateAsset(_propsName string, _starRating uint8, announcer common.Address) (*types.Transaction, error) {
	return _Erc721.Contract.CreateAsset(&_Erc721.TransactOpts, _propsName, _starRating, announcer)
}

// CreateAsset is a paid mutator transaction binding the contract method 0x0324143f.
//
// Solidity: function createAsset(string _propsName, uint8 _starRating, address announcer) returns(uint256)
func (_Erc721 *Erc721TransactorSession) CreateAsset(_propsName string, _starRating uint8, announcer common.Address) (*types.Transaction, error) {
	return _Erc721.Contract.CreateAsset(&_Erc721.TransactOpts, _propsName, _starRating, announcer)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721Transactor) Transfer(opts *bind.TransactOpts, _to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "transfer", _to, _assetId)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721Session) Transfer(_to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Transfer(&_Erc721.TransactOpts, _to, _assetId)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721TransactorSession) Transfer(_to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Transfer(&_Erc721.TransactOpts, _to, _assetId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721Transactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "transferFrom", _from, _to, _assetId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721Session) TransferFrom(_from common.Address, _to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, _from, _to, _assetId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _assetId) payable returns()
func (_Erc721 *Erc721TransactorSession) TransferFrom(_from common.Address, _to common.Address, _assetId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, _from, _to, _assetId)
}

// Erc721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc721 contract.
type Erc721ApprovalIterator struct {
	Event *Erc721Approval // Event containing the contract specifics and raw log

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
func (it *Erc721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Approval)
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
		it.Event = new(Erc721Approval)
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
func (it *Erc721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Approval represents a Approval event raised by the Erc721 contract.
type Erc721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) FilterApproval(opts *bind.FilterOpts, _owner []common.Address, _approved []common.Address, _tokenId []*big.Int) (*Erc721ApprovalIterator, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _approvedRule []interface{}
	for _, _approvedItem := range _approved {
		_approvedRule = append(_approvedRule, _approvedItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Approval", _ownerRule, _approvedRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721ApprovalIterator{contract: _Erc721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc721Approval, _owner []common.Address, _approved []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _ownerRule []interface{}
	for _, _ownerItem := range _owner {
		_ownerRule = append(_ownerRule, _ownerItem)
	}
	var _approvedRule []interface{}
	for _, _approvedItem := range _approved {
		_approvedRule = append(_approvedRule, _approvedItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Approval", _ownerRule, _approvedRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Approval)
				if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed _owner, address indexed _approved, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) ParseApproval(log types.Log) (*Erc721Approval, error) {
	event := new(Erc721Approval)
	if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721CreateAssetIterator is returned from FilterCreateAsset and is used to iterate over the raw logs and unpacked data for CreateAsset events raised by the Erc721 contract.
type Erc721CreateAssetIterator struct {
	Event *Erc721CreateAsset // Event containing the contract specifics and raw log

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
func (it *Erc721CreateAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721CreateAsset)
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
		it.Event = new(Erc721CreateAsset)
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
func (it *Erc721CreateAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721CreateAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721CreateAsset represents a CreateAsset event raised by the Erc721 contract.
type Erc721CreateAsset struct {
	Announcer  common.Address
	TokenId    *big.Int
	PropsName  string
	StarRating uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCreateAsset is a free log retrieval operation binding the contract event 0x50484cca76d2106d636492405f185725c2cafdbe8855cb341deb4b882898a790.
//
// Solidity: event CreateAsset(address announcer, uint256 _tokenId, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721Filterer) FilterCreateAsset(opts *bind.FilterOpts) (*Erc721CreateAssetIterator, error) {

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "CreateAsset")
	if err != nil {
		return nil, err
	}
	return &Erc721CreateAssetIterator{contract: _Erc721.contract, event: "CreateAsset", logs: logs, sub: sub}, nil
}

// WatchCreateAsset is a free log subscription operation binding the contract event 0x50484cca76d2106d636492405f185725c2cafdbe8855cb341deb4b882898a790.
//
// Solidity: event CreateAsset(address announcer, uint256 _tokenId, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721Filterer) WatchCreateAsset(opts *bind.WatchOpts, sink chan<- *Erc721CreateAsset) (event.Subscription, error) {

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "CreateAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721CreateAsset)
				if err := _Erc721.contract.UnpackLog(event, "CreateAsset", log); err != nil {
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

// ParseCreateAsset is a log parse operation binding the contract event 0x50484cca76d2106d636492405f185725c2cafdbe8855cb341deb4b882898a790.
//
// Solidity: event CreateAsset(address announcer, uint256 _tokenId, string _propsName, uint8 _starRating)
func (_Erc721 *Erc721Filterer) ParseCreateAsset(log types.Log) (*Erc721CreateAsset, error) {
	event := new(Erc721CreateAsset)
	if err := _Erc721.contract.UnpackLog(event, "CreateAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc721 contract.
type Erc721TransferIterator struct {
	Event *Erc721Transfer // Event containing the contract specifics and raw log

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
func (it *Erc721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Transfer)
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
		it.Event = new(Erc721Transfer)
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
func (it *Erc721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Transfer represents a Transfer event raised by the Erc721 contract.
type Erc721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) FilterTransfer(opts *bind.FilterOpts, _from []common.Address, _to []common.Address, _tokenId []*big.Int) (*Erc721TransferIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Transfer", _fromRule, _toRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721TransferIterator{contract: _Erc721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc721Transfer, _from []common.Address, _to []common.Address, _tokenId []*big.Int) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _tokenIdRule []interface{}
	for _, _tokenIdItem := range _tokenId {
		_tokenIdRule = append(_tokenIdRule, _tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Transfer", _fromRule, _toRule, _tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Transfer)
				if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed _from, address indexed _to, uint256 indexed _tokenId)
func (_Erc721 *Erc721Filterer) ParseTransfer(log types.Log) (*Erc721Transfer, error) {
	event := new(Erc721Transfer)
	if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
