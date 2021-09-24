package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
)

type eth struct {
	ethClient *ethclient.Client
}

func NewEth(url, appid string) (eth, error) {
	client, err := ethclient.Dial(url + appid)
	if err != nil {
		log.GetLogger().Error("[NewEth] failed", zap.Error(err))
		return eth{}, err
	}
	return eth{ethClient: client}, err
}

func (e eth) Close() {
	e.ethClient.Close()
}
