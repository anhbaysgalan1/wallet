package main

import (
	"fmt"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin"
	"git.huoys.com/middle-end/kratos/pkg/conf/paladin/apollo"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"tp_wallet/config"
	"tp_wallet/internal/block_chain/chain/scan/common"
	"tp_wallet/internal/block_chain/chain/scan/operation"
	"tp_wallet/internal/block_chain/chain/scan/orm"
	"tp_wallet/pkg/log"
)

var wg sync.WaitGroup

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	taskScanBsc()
	wg.Wait()
}

func loadConfig() error {
	err := paladin.Init(apollo.PaladinDriverApollo)
	if err != nil {
		return err
	}

	// init kafka
	kaf, _ := config.NewConfigBlockKafka()
	operation.Producer, err = operation.NewProducer(kaf.KafkaAddr)
	if err != nil {
		return err
	}
	fmt.Println(kaf.TopicCurrencyTransaction)
	fmt.Println(kaf.TopicNftTransaction)
	fmt.Println(kaf.TopicNftCreate)
	operation.TopicH20ForScan = kaf.TopicCurrencyTransaction
	operation.TopicNftForScan = kaf.TopicNftTransaction
	operation.TopicCreateNftForScan = kaf.TopicNftCreate

	// init config
	config.NewBlockBusiness()
	fmt.Println(config.BlockBusiness.RowingBoatContractAddress)
	fmt.Println(config.BlockBusiness.RacerContractAddress)
	fmt.Println(config.BlockBusiness.FFCoinContractAddress)
	fmt.Println(config.BlockBusiness.F1CoinContractAddress)
	fmt.Println(config.BlockBusiness.BNBRechargeAddress)
	fmt.Println(config.BlockBusiness.BNBWithdrawAddress)
	fmt.Println(config.BlockBusiness.KeyScan)
	fmt.Println(config.BlockBusiness.BscBlockNumber)
	fmt.Println(config.BlockBusiness.ChainNetUrl)
	fmt.Println(config.BlockBusiness.MaterialContractAddress)

	//common.H2OContractAddress = config.BlockBusiness.H2OContractAddress
	common.RacingBoatContractAddress = config.BlockBusiness.RowingBoatContractAddress
	common.RacerContractAddress = config.BlockBusiness.RacerContractAddress
	common.FFCoinContractAddress = config.BlockBusiness.FFCoinContractAddress
	common.F1CoinContractAddress = config.BlockBusiness.F1CoinContractAddress
	common.BNBRechargeAddress = config.BlockBusiness.BNBRechargeAddress
	common.BNBWithdrawAddress = config.BlockBusiness.BNBWithdrawAddress
	common.ScanKey = config.BlockBusiness.KeyScan
	common.BeginBlockNumber = config.BlockBusiness.BscBlockNumber
	common.ChainNetUrl = config.BlockBusiness.ChainNetUrl
	common.MaterialContractAddress = config.BlockBusiness.MaterialContractAddress

	// init mongo
	monConf, err := config.ConfNewDB()
	if err != nil {
		return err
	}

	orm.MonCli, err = orm.NewMongoClient(monConf.User, monConf.Password, monConf.Addr)
	if err != nil {
		return err
	}

	_, err = orm.GetHeightByMongo(orm.CollectionBsc)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			//首次没有数据，新插入
			var mh orm.MongoBlockHeight
			mh.NetWork = orm.CollectionBsc
			mh.Height = common.BeginBlockNumber
			errS := orm.SetHeightByMongo(mh)
			if errS != nil {
				log.GetLogger().Error("SetHeightByMongo begin", zap.Error(err))
				return errS
			}
		} else {
			return err
		}
	}

	return nil
}

func taskScanBsc() {
	go func() {
		var i uint64
		for true {
			i++
			height, ok, err := operation.ScanBlockForBsc()
			if err != nil {
				if err.Error() == common.GetBlockErrNotFound {
					time.Sleep(3 * time.Second)
					continue
				}
				log.GetLogger().Error("operation.ScanBlockForBsc", zap.Error(err))
				time.Sleep(3 * time.Second)
				continue
			}
			if ok {
				var mh orm.MongoBlockHeight
				hi, err := strconv.ParseUint(height, 0, 0)
				if err != nil {
					log.GetLogger().Error("strconv.ParseUint(height, 0, 0)", zap.Error(err), zap.Any("number", hi))
					continue
				}
				mh.Height = hi
				mh.NetWork = orm.CollectionBsc
				err = orm.UpdateHeightByMongo(mh)
				if err != nil {
					log.GetLogger().Error("UpdateHeightByMongo", zap.Error(err), zap.Any("number", hi))
					continue
				}
				continue
			}
			time.Sleep(3 * time.Second)
			if i > 10000000000 {
				break
			}
		}
		wg.Done()
	}()
}
