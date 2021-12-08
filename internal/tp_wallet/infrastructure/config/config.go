package config

import (
	"errors"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"go.uber.org/zap"
	"tp_wallet/internal/common"
	"tp_wallet/pkg/database/mongo"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/log"
)

var WalletBusiness *common.ConfigWalletBusiness
var BlockBusiness *common.ConfigBlockBusiness
var Fee *common.ConfigFee
var SysAccountMap map[string]common.AddressForGame
var SysAccountMapByUid map[uint64]common.AddressForGame

var BlockSysAddrToUid map[string]uint64
var CurrencyMap map[string]common.Currency
var CurrencyList []common.Currency

func ConfRedis() (*redis.Config, error) {
	var (
		cfg redis.Config
		err error
	)
	if err = paladin.Get("redis.yaml").UnmarshalYAML(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ConfNewDB() (*mongo.Config, error) {
	var (
		cfg mongo.Config
		err error
	)
	if err = paladin.Get("mongo.yaml").UnmarshalYAML(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewConfigBlockKafka() (*common.ConfigTransferKafka, error) {
	var configBlockKafka common.ConfigTransferKafka
	var err error
	if err = paladin.Get("block_kafka.yaml").UnmarshalYAML(&configBlockKafka); err != nil {
		return nil, err
	}
	return &configBlockKafka, nil
}

func NewWalletBusiness() {
	var err error
	if err = paladin.Get("wallet_business.yaml").UnmarshalYAML(&WalletBusiness); err != nil {
		log.GetLogger().Error("[NewWalletBusiness] paladin UnmarshalYAML failed", zap.Error(err))
		panic(err)
	}
	if len(WalletBusiness.HttpAddr) == 0 || len(WalletBusiness.GrpcAddr) == 0 || len(WalletBusiness.CronAsyncTransfer) == 0 {
		log.GetLogger().Error("[NewWalletBusiness] Parameters failed", zap.Any("config", WalletBusiness))
		panic(errors.New("config init failed"))
	}
}

func NewBlockBusiness() {
	var err error
	if err = paladin.Get("block_chain_business.yaml").UnmarshalYAML(&BlockBusiness); err != nil {
		log.GetLogger().Error("[NewBlockBusiness] paladin UnmarshalYAML failed", zap.Error(err))
		panic(err)
	}
	if len(BlockBusiness.HttpAddr) == 0 || BlockBusiness.H2OSysUid == BlockBusiness.NftSysUid ||
		len(BlockBusiness.H2OUidAddrIncome) == 0 || len(BlockBusiness.H2OUidAddrExpenditure) == 0 ||
		len(BlockBusiness.NftAddrIncome) == 0 || len(BlockBusiness.NftAddrExpenditure) == 0 {
		log.GetLogger().Error("[NewBlockBusiness] Parameters failed", zap.Any("config", BlockBusiness))
		panic(errors.New("config init failed"))
	}
	BlockSysAddrToUid = make(map[string]uint64)
	SysAccountMap = make(map[string]common.AddressForGame)
	SysAccountMapByUid = make(map[uint64]common.AddressForGame)
	for _, item := range BlockBusiness.H2OUidAddrExpenditure {
		BlockSysAddrToUid[item] = BlockBusiness.H2OSysUid
	}
	BlockSysAddrToUid[BlockBusiness.H2OUidAddrIncome] = BlockBusiness.H2OSysUid
	BlockSysAddrToUid[BlockBusiness.H2OContractAddress] = BlockBusiness.H2OSysUid
	for _, item := range BlockBusiness.NftAddrExpenditure {
		BlockSysAddrToUid[item] = BlockBusiness.NftSysUid
	}
	BlockSysAddrToUid[BlockBusiness.NftAddrIncome] = BlockBusiness.NftSysUid
	SysAccountMap[BlockBusiness.H2OSysToken] = common.AddressForGame{
		SysUid:          BlockBusiness.H2OSysUid,
		SysToken:        BlockBusiness.H2OSysToken,
		ContractAddress: BlockBusiness.H2OContractAddress,
		AddrIncome:      BlockBusiness.H2OUidAddrIncome,
		AddrExpenditure: BlockBusiness.H2OUidAddrExpenditure,
	}
	SysAccountMap[BlockBusiness.NftSysToken] = common.AddressForGame{
		SysUid:          BlockBusiness.NftSysUid,
		SysToken:        BlockBusiness.NftSysToken,
		AddrIncome:      BlockBusiness.NftAddrIncome,
		AddrExpenditure: BlockBusiness.NftAddrExpenditure,
	}
	SysAccountMapByUid[BlockBusiness.H2OSysUid] = common.AddressForGame{
		SysUid:          BlockBusiness.H2OSysUid,
		SysToken:        BlockBusiness.H2OSysToken,
		ContractAddress: BlockBusiness.H2OContractAddress,
		AddrIncome:      BlockBusiness.H2OUidAddrIncome,
		AddrExpenditure: BlockBusiness.H2OUidAddrExpenditure,
	}
	SysAccountMapByUid[BlockBusiness.NftSysUid] = common.AddressForGame{
		SysUid:          BlockBusiness.NftSysUid,
		SysToken:        BlockBusiness.NftSysToken,
		AddrIncome:      BlockBusiness.NftAddrIncome,
		AddrExpenditure: BlockBusiness.NftAddrExpenditure,
	}

	for _, items := range BlockBusiness.AddressForGame {
		for _, item := range items.AddrExpenditure {
			BlockSysAddrToUid[item] = items.SysUid
		}
		BlockSysAddrToUid[items.AddrIncome] = items.SysUid
		BlockSysAddrToUid[items.AddrCreate] = items.SysUid
		BlockSysAddrToUid[items.ContractAddress] = items.SysUid
		SysAccountMap[items.SysToken] = items
		SysAccountMapByUid[items.SysUid] = items
	}
}

func NewConfigFee() {
	var err error
	if err = paladin.Get("block_fee.yaml").UnmarshalYAML(&Fee); err != nil {
		log.GetLogger().Error("[NewConfigFee] paladin UnmarshalYAML failed", zap.Error(err))
		panic(err)
	}
	if Fee.GasPrice <= 0 || Fee.GasLimit <= 0 {
		panic("fee gas price and gas limit equal to 0")
	}
}

func NewConfigCurrency() {
	var err error
	CurrencyMap = make(map[string]common.Currency)
	CurrencyList = make([]common.Currency, 0)
	if err = paladin.Get("currency.yaml").UnmarshalYAML(&CurrencyList); err != nil {
		log.GetLogger().Error("[NewConfigCurrency] paladin UnmarshalYAML failed", zap.Error(err))
		panic(err)
	}
	if len(CurrencyList) == 0 {
		panic("currency list length zero")
	}
	for _, item := range CurrencyList {
		CurrencyMap[item.Name] = item
	}
}
