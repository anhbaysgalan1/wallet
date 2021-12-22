package transfer

import "math/big"

// InputForTransfer 转账请求参数
type InputForTransfer struct {
	FromAddress     string   `json:"from_address"`     //提币总账户地址
	ContractAddress string   `json:"contract_address"` //代币合约地址
	ToAddress       string   `json:"to_address"`       //用户第三方接收钱包地址
	Amount          *big.Int `json:"amount"`           //提取数量
	GasLimit        uint64   `json:"gas_limit"`        //gas数量，建议最少200000
	GasPrice        uint64   `json:"gas_price"`        //gas单价，起始单价都是5 GWei
	Nonce           uint64   `json:"nonce"`            //提币总账户地址nonce值
	Private         string   `json:"private"`          //from地址私钥
}

// InputForBNBTransfer 转账请求参数
type InputForBNBTransfer struct {
	FromAddress string   `json:"from_address"` //提币总账户地址
	ToAddress   string   `json:"to_address"`   //用户第三方接收钱包地址
	Amount      *big.Int `json:"amount"`       //提取数量
	GasLimit    uint64   `json:"gas_limit"`    //gas数量，建议最少200000
	GasPrice    uint64   `json:"gas_price"`    //gas单价，起始单价都是5 GWei
	Nonce       uint64   `json:"nonce"`        //提币总账户地址nonce值
	Private     string   `json:"private"`      //from地址私钥
}

type JsonResult struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}
