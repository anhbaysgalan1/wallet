package erc1155

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

// MateralABI is the input ABI used to generate the binding from.
const MateralABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"announcer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"materialName\",\"type\":\"string\"}],\"name\":\"createAsset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"announcer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"materialName\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CreateAsset\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"expand\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contractName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getAsset\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MateralBin is the compiled bytecode used for deploying new contracts.
var MateralBin = "0x60c0604052600c60808190526b118c48120c93c814d411515160a21b60a09081526200002f916001919062000088565b506040805180820190915260088082526713585d195c9a585b60c21b6020909201918252620000619160029162000088565b503480156200006f57600080fd5b50600080546001600160a01b031916331790556200016b565b82805462000096906200012e565b90600052602060002090601f016020900481019282620000ba576000855562000105565b82601f10620000d557805160ff191683800117855562000105565b8280016001018555821562000105579182015b8281111562000105578251825591602001919060010190620000e8565b506200011392915062000117565b5090565b5b8082111562000113576000815560010162000118565b600181811c908216806200014357607f821691505b602082108114156200016557634e487b7160e01b600052602260045260246000fd5b50919050565b6112d8806200017b6000396000f3fe608060405234801561001057600080fd5b50600436106100a85760003560e01c8063a22cb46511610071578063a22cb46514610130578063b63a2c1614610143578063d721fe0214610156578063e985e9c51461015e578063eac8f5b8146101aa578063fba0ee64146101bd57600080fd5b8062fdd58e146100ad5780630febdd49146100d35780634e1273f4146100e857806375d0c0dc14610108578063874b3afc1461011d575b600080fd5b6100c06100bb366004610e01565b6101d0565b6040519081526020015b60405180910390f35b6100e66100e1366004610d83565b610269565b005b6100fb6100f6366004610e2b565b6102f5565b6040516100ca919061109c565b61011061041f565b6040516100ca91906110e4565b6100e661012b366004610fba565b6104ad565b6100e661013e366004610dc5565b610584565b6100c0610151366004610f0b565b610593565b6101106106db565b61019a61016c366004610ccb565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205460ff1690565b60405190151581526020016100ca565b6100c06101b8366004610ef2565b6106e8565b6100e66101cb366004610cfe565b61075f565b60006001600160a01b0383166102415760405162461bcd60e51b815260206004820152602b60248201527f455243313135353a2062616c616e636520717565727920666f7220746865207a60448201526a65726f206164647265737360a81b60648201526084015b60405180910390fd5b5060009081526004602090815260408083206001600160a01b03949094168352929052205490565b6001600160a01b0384163314806102855750610285843361016c565b6102e35760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2063616c6c6572206973206e6f74206f776e6572206e6f7260448201526808185c1c1c9bdd995960ba1b6064820152608401610238565b6102ef848484846107ee565b50505050565b6060815183511461035a5760405162461bcd60e51b815260206004820152602960248201527f455243313135353a206163636f756e747320616e6420696473206c656e677468604482015268040dad2e6dac2e8c6d60bb1b6064820152608401610238565b6000835167ffffffffffffffff8111156103765761037661128c565b60405190808252806020026020018201604052801561039f578160200160208202803683370190505b50905060005b8451811015610417576103ea8582815181106103c3576103c3611276565b60200260200101518583815181106103dd576103dd611276565b60200260200101516101d0565b8282815181106103fc576103fc611276565b602090810291909101015261041081611245565b90506103a5565b509392505050565b6002805461042c9061120a565b80601f01602080910402602001604051908101604052809291908181526020018280546104589061120a565b80156104a55780601f1061047a576101008083540402835291602001916104a5565b820191906000526020600020905b81548152906001019060200180831161048857829003601f168201915b505050505081565b60035482106104f55760405162461bcd60e51b8152602060048201526014602482015273082e6e6cae874d2c840d2e640e8dede40d0d2ced60631b6044820152606401610238565b60008281526004602090815260408083203384529091529020548181101561052f5760405162461bcd60e51b81526004016102389061113c565b600083815260046020908152604080832033845290915290208282039055600380548391908590811061056457610564611276565b600091825260209091206001600290920201018054919091039055505050565b61058f3383836108f2565b5050565b600080546001600160a01b031633146105ee5760405162461bcd60e51b815260206004820152601e60248201527f41737365743a6d73672e73656e646572206e6f74206973206f776e65722100006044820152606401610238565b6040805180820190915282815260208082018690526003805460018101825560009190915282518051849360029093027fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b019261064f928492910190610ba4565b506020919091015160019182015560035460009161066c916111f3565b60008181526004602090815260408083206001600160a01b038a16845290915290819020889055519091507f993a624bd9daf00efa6193dcea760e2784bc9dc6f2559f3e9f9b316e301e19fc906106ca908790849088908b90611064565b60405180910390a195945050505050565b6001805461042c9061120a565b60035460009082106107335760405162461bcd60e51b8152602060048201526014602482015273082e6e6cae874d2c840d2e640e8dede40d0d2ced60631b6044820152606401610238565b6003828154811061074657610746611276565b9060005260206000209060020201600101549050919050565b6001600160a01b03841633148061077b575061077b843361016c565b6107e25760405162461bcd60e51b815260206004820152603260248201527f455243313135353a207472616e736665722063616c6c6572206973206e6f74206044820152711bdddb995c881b9bdc88185c1c1c9bdd995960721b6064820152608401610238565b6102ef848484846109d3565b6001600160a01b0383166108145760405162461bcd60e51b8152600401610238906110f7565b60008281526004602090815260408083206001600160a01b03881684529091529020543390828110156108595760405162461bcd60e51b81526004016102389061113c565b60008481526004602090815260408083206001600160a01b038a81168552925280832086850390559087168252812080548592906108989084906111db565b909155505060408051858152602081018590526001600160a01b038088169289821692918616917fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62910160405180910390a4505050505050565b816001600160a01b0316836001600160a01b031614156109665760405162461bcd60e51b815260206004820152602960248201527f455243313135353a2073657474696e6720617070726f76616c20737461747573604482015268103337b91039b2b63360b91b6064820152608401610238565b6001600160a01b03838116600081815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b8051825114610a355760405162461bcd60e51b815260206004820152602860248201527f455243313135353a2069647320616e6420616d6f756e7473206c656e677468206044820152670dad2e6dac2e8c6d60c31b6064820152608401610238565b6001600160a01b038316610a5b5760405162461bcd60e51b8152600401610238906110f7565b3360005b8351811015610b45576000848281518110610a7c57610a7c611276565b602002602001015190506000848381518110610a9a57610a9a611276565b60209081029190910181015160008481526004835260408082206001600160a01b038d168352909352919091205490915081811015610aeb5760405162461bcd60e51b81526004016102389061113c565b60008381526004602090815260408083206001600160a01b038d8116855292528083208585039055908a16825281208054849290610b2a9084906111db565b9250508190555050505080610b3e90611245565b9050610a5f565b50836001600160a01b0316856001600160a01b0316826001600160a01b03167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb8686604051610b959291906110b6565b60405180910390a45050505050565b828054610bb09061120a565b90600052602060002090601f016020900481019282610bd25760008555610c18565b82601f10610beb57805160ff1916838001178555610c18565b82800160010185558215610c18579182015b82811115610c18578251825591602001919060010190610bfd565b50610c24929150610c28565b5090565b5b80821115610c245760008155600101610c29565b80356001600160a01b0381168114610c5457600080fd5b919050565b600082601f830112610c6a57600080fd5b81356020610c7f610c7a836111b7565b611186565b80838252828201915082860187848660051b8901011115610c9f57600080fd5b60005b85811015610cbe57813584529284019290840190600101610ca2565b5090979650505050505050565b60008060408385031215610cde57600080fd5b610ce783610c3d565b9150610cf560208401610c3d565b90509250929050565b60008060008060808587031215610d1457600080fd5b610d1d85610c3d565b9350610d2b60208601610c3d565b9250604085013567ffffffffffffffff80821115610d4857600080fd5b610d5488838901610c59565b93506060870135915080821115610d6a57600080fd5b50610d7787828801610c59565b91505092959194509250565b60008060008060808587031215610d9957600080fd5b610da285610c3d565b9350610db060208601610c3d565b93969395505050506040820135916060013590565b60008060408385031215610dd857600080fd5b610de183610c3d565b915060208301358015158114610df657600080fd5b809150509250929050565b60008060408385031215610e1457600080fd5b610e1d83610c3d565b946020939093013593505050565b60008060408385031215610e3e57600080fd5b823567ffffffffffffffff80821115610e5657600080fd5b818501915085601f830112610e6a57600080fd5b81356020610e7a610c7a836111b7565b8083825282820191508286018a848660051b8901011115610e9a57600080fd5b600096505b84871015610ec457610eb081610c3d565b835260019690960195918301918301610e9f565b5096505086013592505080821115610edb57600080fd5b50610ee885828601610c59565b9150509250929050565b600060208284031215610f0457600080fd5b5035919050565b600080600060608486031215610f2057600080fd5b833592506020610f31818601610c3d565b9250604085013567ffffffffffffffff80821115610f4e57600080fd5b818701915087601f830112610f6257600080fd5b813581811115610f7457610f7461128c565b610f86601f8201601f19168501611186565b91508082528884828501011115610f9c57600080fd5b80848401858401376000848284010152508093505050509250925092565b60008060408385031215610fcd57600080fd5b50508035926020909101359150565b600081518084526020808501945080840160005b8381101561100c57815187529582019590820190600101610ff0565b509495945050505050565b6000815180845260005b8181101561103d57602081850181015186830182015201611021565b8181111561104f576000602083870101525b50601f01601f19169290920160200192915050565b60018060a01b038516815283602082015260806040820152600061108b6080830185611017565b905082606083015295945050505050565b6020815260006110af6020830184610fdc565b9392505050565b6040815260006110c96040830185610fdc565b82810360208401526110db8185610fdc565b95945050505050565b6020815260006110af6020830184611017565b60208082526025908201527f455243313135353a207472616e7366657220746f20746865207a65726f206164604082015264647265737360d81b606082015260800190565b6020808252602a908201527f455243313135353a20696e73756666696369656e742062616c616e636520666f60408201526939103a3930b739b332b960b11b606082015260800190565b604051601f8201601f1916810167ffffffffffffffff811182821017156111af576111af61128c565b604052919050565b600067ffffffffffffffff8211156111d1576111d161128c565b5060051b60200190565b600082198211156111ee576111ee611260565b500190565b60008282101561120557611205611260565b500390565b600181811c9082168061121e57607f821691505b6020821081141561123f57634e487b7160e01b600052602260045260246000fd5b50919050565b600060001982141561125957611259611260565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea264697066735822122017551bc41a73964d23f7038bcfef5d57d3b8af3f98f876ef453103efef72330e64736f6c63430008070033"

// DeployMateral deploys a new Ethereum contract, binding an instance of Materal to it.
func DeployMateral(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Materal, error) {
	parsed, err := abi.JSON(strings.NewReader(MateralABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MateralBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Materal{MateralCaller: MateralCaller{contract: contract}, MateralTransactor: MateralTransactor{contract: contract}, MateralFilterer: MateralFilterer{contract: contract}}, nil
}

// Materal is an auto generated Go binding around an Ethereum contract.
type Materal struct {
	MateralCaller     // Read-only binding to the contract
	MateralTransactor // Write-only binding to the contract
	MateralFilterer   // Log filterer for contract events
}

// MateralCaller is an auto generated read-only Go binding around an Ethereum contract.
type MateralCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MateralTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MateralTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MateralFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MateralFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MateralSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MateralSession struct {
	Contract     *Materal          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MateralCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MateralCallerSession struct {
	Contract *MateralCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// MateralTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MateralTransactorSession struct {
	Contract     *MateralTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MateralRaw is an auto generated low-level Go binding around an Ethereum contract.
type MateralRaw struct {
	Contract *Materal // Generic contract binding to access the raw methods on
}

// MateralCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MateralCallerRaw struct {
	Contract *MateralCaller // Generic read-only contract binding to access the raw methods on
}

// MateralTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MateralTransactorRaw struct {
	Contract *MateralTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMateral creates a new instance of Materal, bound to a specific deployed contract.
func NewMateral(address common.Address, backend bind.ContractBackend) (*Materal, error) {
	contract, err := bindMateral(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Materal{MateralCaller: MateralCaller{contract: contract}, MateralTransactor: MateralTransactor{contract: contract}, MateralFilterer: MateralFilterer{contract: contract}}, nil
}

// NewMateralCaller creates a new read-only instance of Materal, bound to a specific deployed contract.
func NewMateralCaller(address common.Address, caller bind.ContractCaller) (*MateralCaller, error) {
	contract, err := bindMateral(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MateralCaller{contract: contract}, nil
}

// NewMateralTransactor creates a new write-only instance of Materal, bound to a specific deployed contract.
func NewMateralTransactor(address common.Address, transactor bind.ContractTransactor) (*MateralTransactor, error) {
	contract, err := bindMateral(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MateralTransactor{contract: contract}, nil
}

// NewMateralFilterer creates a new log filterer instance of Materal, bound to a specific deployed contract.
func NewMateralFilterer(address common.Address, filterer bind.ContractFilterer) (*MateralFilterer, error) {
	contract, err := bindMateral(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MateralFilterer{contract: contract}, nil
}

// bindMateral binds a generic wrapper to an already deployed contract.
func bindMateral(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MateralABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Materal *MateralRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Materal.Contract.MateralCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Materal *MateralRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Materal.Contract.MateralTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Materal *MateralRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Materal.Contract.MateralTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Materal *MateralCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Materal.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Materal *MateralTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Materal.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Materal *MateralTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Materal.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Materal *MateralCaller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Materal *MateralSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Materal.Contract.BalanceOf(&_Materal.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Materal *MateralCallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Materal.Contract.BalanceOf(&_Materal.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Materal *MateralCaller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Materal *MateralSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Materal.Contract.BalanceOfBatch(&_Materal.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Materal *MateralCallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Materal.Contract.BalanceOfBatch(&_Materal.CallOpts, accounts, ids)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() view returns(string)
func (_Materal *MateralCaller) ContractName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "contractName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() view returns(string)
func (_Materal *MateralSession) ContractName() (string, error) {
	return _Materal.Contract.ContractName(&_Materal.CallOpts)
}

// ContractName is a free data retrieval call binding the contract method 0x75d0c0dc.
//
// Solidity: function contractName() view returns(string)
func (_Materal *MateralCallerSession) ContractName() (string, error) {
	return _Materal.Contract.ContractName(&_Materal.CallOpts)
}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 id) view returns(uint256)
func (_Materal *MateralCaller) GetAsset(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "getAsset", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 id) view returns(uint256)
func (_Materal *MateralSession) GetAsset(id *big.Int) (*big.Int, error) {
	return _Materal.Contract.GetAsset(&_Materal.CallOpts, id)
}

// GetAsset is a free data retrieval call binding the contract method 0xeac8f5b8.
//
// Solidity: function getAsset(uint256 id) view returns(uint256)
func (_Materal *MateralCallerSession) GetAsset(id *big.Int) (*big.Int, error) {
	return _Materal.Contract.GetAsset(&_Materal.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Materal *MateralCaller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Materal *MateralSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Materal.Contract.IsApprovedForAll(&_Materal.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Materal *MateralCallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Materal.Contract.IsApprovedForAll(&_Materal.CallOpts, account, operator)
}

// PlatformName is a free data retrieval call binding the contract method 0xd721fe02.
//
// Solidity: function platformName() view returns(string)
func (_Materal *MateralCaller) PlatformName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Materal.contract.Call(opts, &out, "platformName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// PlatformName is a free data retrieval call binding the contract method 0xd721fe02.
//
// Solidity: function platformName() view returns(string)
func (_Materal *MateralSession) PlatformName() (string, error) {
	return _Materal.Contract.PlatformName(&_Materal.CallOpts)
}

// PlatformName is a free data retrieval call binding the contract method 0xd721fe02.
//
// Solidity: function platformName() view returns(string)
func (_Materal *MateralCallerSession) PlatformName() (string, error) {
	return _Materal.Contract.PlatformName(&_Materal.CallOpts)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xb63a2c16.
//
// Solidity: function createAsset(uint256 amount, address announcer, string materialName) returns(uint256)
func (_Materal *MateralTransactor) CreateAsset(opts *bind.TransactOpts, amount *big.Int, announcer common.Address, materialName string) (*types.Transaction, error) {
	return _Materal.contract.Transact(opts, "createAsset", amount, announcer, materialName)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xb63a2c16.
//
// Solidity: function createAsset(uint256 amount, address announcer, string materialName) returns(uint256)
func (_Materal *MateralSession) CreateAsset(amount *big.Int, announcer common.Address, materialName string) (*types.Transaction, error) {
	return _Materal.Contract.CreateAsset(&_Materal.TransactOpts, amount, announcer, materialName)
}

// CreateAsset is a paid mutator transaction binding the contract method 0xb63a2c16.
//
// Solidity: function createAsset(uint256 amount, address announcer, string materialName) returns(uint256)
func (_Materal *MateralTransactorSession) CreateAsset(amount *big.Int, announcer common.Address, materialName string) (*types.Transaction, error) {
	return _Materal.Contract.CreateAsset(&_Materal.TransactOpts, amount, announcer, materialName)
}

// Expand is a paid mutator transaction binding the contract method 0x874b3afc.
//
// Solidity: function expand(uint256 id, uint256 amount) returns()
func (_Materal *MateralTransactor) Expand(opts *bind.TransactOpts, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.contract.Transact(opts, "expand", id, amount)
}

// Expand is a paid mutator transaction binding the contract method 0x874b3afc.
//
// Solidity: function expand(uint256 id, uint256 amount) returns()
func (_Materal *MateralSession) Expand(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.Contract.Expand(&_Materal.TransactOpts, id, amount)
}

// Expand is a paid mutator transaction binding the contract method 0x874b3afc.
//
// Solidity: function expand(uint256 id, uint256 amount) returns()
func (_Materal *MateralTransactorSession) Expand(id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.Contract.Expand(&_Materal.TransactOpts, id, amount)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts) returns()
func (_Materal *MateralTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _Materal.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, amounts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts) returns()
func (_Materal *MateralSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _Materal.Contract.SafeBatchTransferFrom(&_Materal.TransactOpts, from, to, ids, amounts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0xfba0ee64.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts) returns()
func (_Materal *MateralTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _Materal.Contract.SafeBatchTransferFrom(&_Materal.TransactOpts, from, to, ids, amounts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount) returns()
func (_Materal *MateralTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.contract.Transact(opts, "safeTransferFrom", from, to, id, amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount) returns()
func (_Materal *MateralSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.Contract.SafeTransferFrom(&_Materal.TransactOpts, from, to, id, amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x0febdd49.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount) returns()
func (_Materal *MateralTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Materal.Contract.SafeTransferFrom(&_Materal.TransactOpts, from, to, id, amount)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Materal *MateralTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Materal.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Materal *MateralSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Materal.Contract.SetApprovalForAll(&_Materal.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Materal *MateralTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Materal.Contract.SetApprovalForAll(&_Materal.TransactOpts, operator, approved)
}

// MateralApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Materal contract.
type MateralApprovalForAllIterator struct {
	Event *MateralApprovalForAll // Event containing the contract specifics and raw log

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
func (it *MateralApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MateralApprovalForAll)
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
		it.Event = new(MateralApprovalForAll)
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
func (it *MateralApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MateralApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MateralApprovalForAll represents a ApprovalForAll event raised by the Materal contract.
type MateralApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Materal *MateralFilterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*MateralApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Materal.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &MateralApprovalForAllIterator{contract: _Materal.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Materal *MateralFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *MateralApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Materal.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MateralApprovalForAll)
				if err := _Materal.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Materal *MateralFilterer) ParseApprovalForAll(log types.Log) (*MateralApprovalForAll, error) {
	event := new(MateralApprovalForAll)
	if err := _Materal.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MateralCreateAssetIterator is returned from FilterCreateAsset and is used to iterate over the raw logs and unpacked data for CreateAsset events raised by the Materal contract.
type MateralCreateAssetIterator struct {
	Event *MateralCreateAsset // Event containing the contract specifics and raw log

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
func (it *MateralCreateAssetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MateralCreateAsset)
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
		it.Event = new(MateralCreateAsset)
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
func (it *MateralCreateAssetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MateralCreateAssetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MateralCreateAsset represents a CreateAsset event raised by the Materal contract.
type MateralCreateAsset struct {
	Announcer    common.Address
	Id           *big.Int
	MaterialName string
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterCreateAsset is a free log retrieval operation binding the contract event 0x993a624bd9daf00efa6193dcea760e2784bc9dc6f2559f3e9f9b316e301e19fc.
//
// Solidity: event CreateAsset(address announcer, uint256 id, string materialName, uint256 amount)
func (_Materal *MateralFilterer) FilterCreateAsset(opts *bind.FilterOpts) (*MateralCreateAssetIterator, error) {

	logs, sub, err := _Materal.contract.FilterLogs(opts, "CreateAsset")
	if err != nil {
		return nil, err
	}
	return &MateralCreateAssetIterator{contract: _Materal.contract, event: "CreateAsset", logs: logs, sub: sub}, nil
}

// WatchCreateAsset is a free log subscription operation binding the contract event 0x993a624bd9daf00efa6193dcea760e2784bc9dc6f2559f3e9f9b316e301e19fc.
//
// Solidity: event CreateAsset(address announcer, uint256 id, string materialName, uint256 amount)
func (_Materal *MateralFilterer) WatchCreateAsset(opts *bind.WatchOpts, sink chan<- *MateralCreateAsset) (event.Subscription, error) {

	logs, sub, err := _Materal.contract.WatchLogs(opts, "CreateAsset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MateralCreateAsset)
				if err := _Materal.contract.UnpackLog(event, "CreateAsset", log); err != nil {
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

// ParseCreateAsset is a log parse operation binding the contract event 0x993a624bd9daf00efa6193dcea760e2784bc9dc6f2559f3e9f9b316e301e19fc.
//
// Solidity: event CreateAsset(address announcer, uint256 id, string materialName, uint256 amount)
func (_Materal *MateralFilterer) ParseCreateAsset(log types.Log) (*MateralCreateAsset, error) {
	event := new(MateralCreateAsset)
	if err := _Materal.contract.UnpackLog(event, "CreateAsset", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MateralTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Materal contract.
type MateralTransferBatchIterator struct {
	Event *MateralTransferBatch // Event containing the contract specifics and raw log

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
func (it *MateralTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MateralTransferBatch)
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
		it.Event = new(MateralTransferBatch)
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
func (it *MateralTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MateralTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MateralTransferBatch represents a TransferBatch event raised by the Materal contract.
type MateralTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Materal *MateralFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*MateralTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Materal.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MateralTransferBatchIterator{contract: _Materal.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Materal *MateralFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *MateralTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Materal.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MateralTransferBatch)
				if err := _Materal.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Materal *MateralFilterer) ParseTransferBatch(log types.Log) (*MateralTransferBatch, error) {
	event := new(MateralTransferBatch)
	if err := _Materal.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MateralTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Materal contract.
type MateralTransferSingleIterator struct {
	Event *MateralTransferSingle // Event containing the contract specifics and raw log

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
func (it *MateralTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MateralTransferSingle)
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
		it.Event = new(MateralTransferSingle)
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
func (it *MateralTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MateralTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MateralTransferSingle represents a TransferSingle event raised by the Materal contract.
type MateralTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Materal *MateralFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*MateralTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Materal.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MateralTransferSingleIterator{contract: _Materal.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Materal *MateralFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *MateralTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Materal.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MateralTransferSingle)
				if err := _Materal.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Materal *MateralFilterer) ParseTransferSingle(log types.Log) (*MateralTransferSingle, error) {
	event := new(MateralTransferSingle)
	if err := _Materal.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
