package common

// for contract address
var (
	H2OContractAddress        string //H2O代币合约地址
	RacingBoatContractAddress string //赛艇NFT合约地址
	RacerContractAddress      string //赛手NFT合约地址
	FFCoinContractAddress     string //FF代币合约地址
	F1CoinContractAddress     string //F1代币合约地址
	MaterialContractAddress   string //材料NFT合约地址
	BNBRechargeAddress        string //BNB充值总地址
	BNBWithdrawAddress        string //BNB提现总地址
	ScanKey                   string //bsc api key
	BeginBlockNumber          uint64 //bsc 扫块起始高度
	ChainNetUrl               string //扫块服务网络：testNet or mainNet
)

// for 交易类型
const (
	UnknownCode                 = 0  //未知交易
	H2OTransferCode             = 1  //H2O转账交易
	H2OApprovalCode             = 2  //H2O授权交易
	H2OTransferFromCode         = 3  //H2O授权转账交易
	RowingNftCreateCode         = 4  //赛艇NFT创建交易
	RowingNftTransferCode       = 5  //赛艇NFT转让交易
	RowingNftApprovalCode       = 6  //赛艇NFT授权交易
	RowingNftTransferFromCode   = 7  //赛艇NFT授权转让交易
	RacerNftCreateCode          = 8  //赛手NFT创建交易
	RacerNftTransferCode        = 9  //赛手NFT转让交易
	RacerNftApprovalCode        = 10 //赛手NFT授权交易
	RacerNftTransferFromCode    = 11 //赛手NFT授权转让交易
	FFTransferCode              = 12 //FF转账交易
	FFApprovalCode              = 13 //FF授权交易
	FFTransferFromCode          = 14 //FF授权转账交易
	F1TransferCode              = 15 //F1转账交易
	F1ApprovalCode              = 16 //F1授权交易
	F1TransferFromCode          = 17 //F1授权转账交易
	BNBRechargeCode             = 18 //BNB充值交易
	BNBWithdrawCode             = 19 //BNB提现交易
	MaterialNftCreateCode       = 20 //材料NFT创建交易
	MaterialNftTransferCode     = 21 //材料NFT转让交易
	MaterialNftApprovalCode     = 22 //材料NFT授权交易
	MaterialNftTransferFromCode = 23 //材料NFT授权转让交易
	MaterialNftExpandCode       = 24 //材料NFT授权转让交易
)

const (
	FFCurrency  = "ff"
	F1Currency  = "f1"
	H2OCurrency = "h2o"
	BnBCurrency = "bnb"
)

const (
	TransferPackInput       = "0xa9059cbb"
	ApprovalPackInput       = "0x095ea7b3"
	TransferPackFromInput   = "0x23b872dd"
	CreateAsset721PackInput = "0x0324143f"
)

const (
	CreateMaterialInput   = "0xb63a2c16"
	TransferMaterialInput = "0x0febdd49"
	ExpandMaterialInput   = "0x874b3afc"
)

const (
	GetTransactionErr   = "GetTransactionReceiptError"
	GetBlockErrNotFound = "not found"
)
