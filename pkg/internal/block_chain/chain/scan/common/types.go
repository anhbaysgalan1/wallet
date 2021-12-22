package common

/*
	Code =
	UnknownCode               = 0  //未知交易
	H2OTransferCode           = 1  //H2O转账交易
	H2OApprovalCode           = 2  //H2O授权交易
	H2OTransferFromCode       = 3  //H2O授权转账交易
	RowingNftCreateCode       = 4  //赛艇NFT创建交易
	RowingNftTransferCode     = 5  //赛艇NFT转让交易
	RowingNftApprovalCode     = 6  //赛艇NFT授权交易
	RowingNftTransferFromCode = 7  //赛艇NFT授权转让交易
	RacerNftCreateCode        = 8  //赛手NFT创建交易
	RacerNftTransferCode      = 9  //赛手NFT转让交易
	RacerNftApprovalCode      = 10 //赛手NFT授权交易
	RacerNftTransferFromCode  = 11 //赛手NFT授权转让交易
	FFTransferCode            = 12 //FF转账交易
	FFApprovalCode            = 13 //FF授权交易
	FFTransferFromCode        = 14 //FF授权转账交易
	F1TransferCode            = 15 //F1转账交易
	F1ApprovalCode            = 16 //F1授权交易
	F1TransferFromCode        = 17 //F1授权转账交易
	BNBRechargeCode           = 18 //BNB充值交易
	BNBWithdrawCode           = 19 //BNB提现交易
	MaterialNftCreateCode       = 20 //材料NFT创建交易
	MaterialNftTransferCode     = 21 //材料NFT转让交易
	MaterialNftApprovalCode     = 22 //材料NFT授权交易
	MaterialNftTransferFromCode = 23 //材料NFT授权转让交易
*/
type PushInput struct {
	Code uint   `bson:"code"json:"code"`
	Data []byte `bson:"data"json:"data"`
}

type BNBTransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Amount      string `bson:"amount"json:"amount"`             //交易BNB的数量
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址nonce值,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高 数值是16进制字符串：例如block_number= 0xac0977
	Currency    string `bson:"currency"json:"currency"`         //币种名称
}

type FFTransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Amount      string `bson:"amount"json:"amount"`             //交易FF的数量
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址nonce值,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高 数值是16进制字符串：例如block_number= 0xac0977
	Currency    string `bson:"currency"json:"currency"`         //币种名称
}

type F1TransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Amount      string `bson:"amount"json:"amount"`             //交易F1的数量
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址nonce值,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高 数值是16进制字符串：例如block_number= 0xac0977
	Currency    string `bson:"currency"json:"currency"`         //币种名称
}

type H20TransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Amount      string `bson:"amount"json:"amount"`             //交易H2O的数量
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址nonce值,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高 数值是16进制字符串：例如block_number= 0xac0977
	Currency    string `bson:"currency"json:"currency"`         //币种名称
}

type Nft721TransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址Nonce,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	NftToken    string `bson:"nft_token"json:"nft_token"`       //NftToken
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高,数值是16进制字符串：例如block_number= 0xac0977
}

type Nft721CreateForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址Nonce,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	NftToken    string `bson:"nft_token"json:"nft_token"`       //NftToken
	PropsName   string `bson:"props_name"json:"props_name"`     //道具名称
	StarRating  string `bson:"star_rating"json:"star_rating"`   //星级
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高,数值是16进制字符串：例如block_number= 0xac0977
}

type Nft1155TransactionForPush struct {
	From        string `bson:"from"json:"from"`                 //钱包地址
	To          string `bson:"to"json:"to"`                     //钱包地址
	Contract    string `bson:"contract"json:"contract"`         //合约地址
	Nonce       string `bson:"nonce"json:"nonce"`               //from地址Nonce,数值是16进制字符串：例如nonce= 0x59
	Hash        string `bson:"hash"json:"hash"`                 //交易hash
	NftToken    string `bson:"nft_token"json:"nft_token"`       //材料NftTokenID
	Amount      string `bson:"amount"json:"amount"`             //材料nft的数量
	Status      string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber string `bson:"block_number"json:"block_number"` //块高,数值是16进制字符串：例如block_number= 0xac0977
}

type Nft1155CreateForPush struct {
	From         string `bson:"from"json:"from"`                 //钱包地址
	To           string `bson:"to"json:"to"`                     //钱包地址
	Contract     string `bson:"contract"json:"contract"`         //合约地址
	Nonce        string `bson:"nonce"json:"nonce"`               //from地址Nonce,数值是16进制字符串：例如nonce= 0x59
	Hash         string `bson:"hash"json:"hash"`                 //交易hash
	NftToken     string `bson:"nft"json:"nft"`                   //材料NftTokenID
	Amount       string `bson:"amount"json:"amount"`             //材料nft的数量
	MaterialName string `bson:"props_name"json:"props_name"`     //材料名称
	Status       string `bson:"status"json:"status"`             //交易状态success or failed
	BlockNumber  string `bson:"block_number"json:"block_number"` //块高,数值是16进制字符串：例如block_number= 0xac0977
}
