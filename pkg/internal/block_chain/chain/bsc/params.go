package bsc

// ethereum mainnet
const (
	EthGasPrice           = "module=proxy&action=eth_gasPrice"
	EthGasTracker         = "module=gastracker&action=gasoracle"
	EthBlock              = "module=proxy&action=eth_getBlockByNumber"
	EthTransactionReceipt = "module=proxy&action=eth_getTransactionReceipt"
	EthTransactionCount   = "module=proxy&action=eth_getTransactionCount"
	EthBalance            = "module=account&action=balance"
	EthSendRawTransaction = "module=proxy&action=eth_sendRawTransaction"
	EthCallContract       = "module=proxy&action=eth_call"
	EthNodeCount          = "module=stats&action=nodecount"
	ContractForErcUSDT    = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	ContractForBscUSDT    = "0x55d398326f99059ff775485246999027b3197955"
	PayAddress            = "0xd75596573b4e691e2ee7cb3b5618b8ab8618c7d5"
	//UrlMainNet            = "https://api.bscscan.com/api?"
	//UrlTestNet            = "https://api-testnet.bscscan.com/api?"
)

const (
	TestKey = "RYD212M9HV1MNA8WT7QPGNZ1I5UYXFIARI"
)
