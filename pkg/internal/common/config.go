package common

import "tp_wallet/pkg/tool"

type ConfigTransferKafka struct {
	KafkaAddr                []string `yaml:"KafkaAddr"`
	TopicCurrencyTransaction string   `yaml:"TopicCurrencyTransaction"`
	GroupCurrencyTransaction string   `yaml:"GroupCurrencyTransaction"`
	TopicNftTransaction      string   `yaml:"TopicNftTransaction"`
	GroupNftTransaction      string   `yaml:"GroupNftTransaction"`
	TopicNftCreate           string   `yaml:"TopicNftCreate"`
	GroupNftCreate           string   `yaml:"GroupNftCreate"`
	TopicProps               string   `yaml:"TopicProps"`
	GroupProps               string   `yaml:"GroupProps"`
}

type ConfigWalletBusiness struct {
	HttpName          string `yaml:"httpName"` // 服务名称
	HttpAddr          string `yaml:"httpAddr"`
	GrpcAddr          string `yaml:"grpcAddr"`
	ReadTimeout       int    `yaml:"readTimeout"`  // 单位s
	WriteTimeout      int    `yaml:"writeTimeout"` // 单位s
	CronAsyncTransfer string `yaml:"cronAsyncTransfer"`
}

type ConfigBlockBusiness struct {
	HttpAddr       string `yaml:"HttpAddr"`       //服务部署地址
	KeyScan        string `yaml:"KeyScan"`        // 用于区块链查询key
	KeyTransfer    string `yaml:"KeyTransfer"`    // 用于区块链查询key
	BscBlockNumber uint64 `yaml:"BscBlockNumber"` //扫块服务起始块高，后面会记录到mongo
	ChainNetUrl    string `yaml:"ChainNetUrl"`    //扫块服务网络：testNet or mainNet
	NetType        string `yaml:"NetType"`        //交易服务网络类型：testNet or mainNet
	//H2OContractAddress        string          `yaml:"H2OContractAddress"`        //H2O代币地址
	MaterialContractAddress   string          `yaml:"MaterialContractAddress"`   //材料NFT合约地址
	RowingBoatContractAddress string          `yaml:"RowingBoatContractAddress"` //赛艇合约地址
	RacerContractAddress      string          `yaml:"RacerContractAddress"`      //赛手合约地址
	FFCoinContractAddress     string          `yaml:"FFCoinContractAddress"`     //FF代币合约地址
	F1CoinContractAddress     string          `yaml:"F1CoinContractAddress"`     //F1代币合约地址
	BNBRechargeAddress        string          `yaml:"BNBRechargeAddress"`        //BNB充值总地址
	BNBWithdrawAddress        string          `yaml:"BNBWithdrawAddress"`        //BNB提现总地址
	AddressForSys             []AddressForSys `yaml:"AddressForSys"`             // 多个游戏收入支出地址
}

type AddressForSys struct {
	SysUid          uint64   `yaml:"SysUid"`          // 系统用户id
	SysToken        string   `yaml:"SysToken"`        // 交易标识
	ContractAddress string   `yaml:"ContractAddress"` // 合约地址
	AddrIncome      string   `yaml:"AddrIncome"`      // 交易收入地址
	AddrCreate      string   `yaml:"AddrCreate"`      // 创建nft地址
	AddrExpenditure []string `yaml:"AddrExpenditure"` // 交易支出地址
}

type ConfigFee struct { // 手续费
	GasLimit uint64 `yaml:"GasLimit"`
	GasPrice uint64 `yaml:"GasPrice"`
}

type Currency struct {
	Name       string         `yaml:"Name"`       // 币种名称
	Status     tool.Switch    `yaml:"Status"`     // 状态开关  1：转账，2：提现，3：充值
	MinCash    string         `yaml:"MinCash"`    // 最小提现金额
	DefaultNet string         `yaml:"DefaultNet"` // 默认网络
	Net        map[string]Net `yaml:"Net"`        // 网络
}

type Net struct {
	NetName     string `yaml:"NetName"`     // 网络名称
	NetType     string `yaml:"NetType"`     // 网络类型
	GasCurrency string `yaml:"GasCurrency"` // 手续费币种
	Gas         string `yaml:"Gas"`         // 手续费
	GasPrice    string `yaml:"GasPrice"`
	GasLimit    string `yaml:"GasLimit"`
	Weight      int64  `yaml:"Weight"` // 排序权重
}

func (c Currency) GetTransferStatus() bool {
	return c.Status.CheckTurnOn(1)
}

func (c Currency) GetCashStatus() bool {
	return c.Status.CheckTurnOn(2)
}

func (c Currency) GetChargeStatus() bool {
	return c.Status.CheckTurnOn(3)
}
