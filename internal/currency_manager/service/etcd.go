package service

import (
	"context"
	"github.com/leaf-rain/wallet/internal/currency_manager/entity"
	"github.com/leaf-rain/wallet/pkg/etcd_config"
	"time"
)

func NewCurrencyCfg() *entity.CurrencyList {
	cfg := &entity.CurrencyList{}
	ecf, err := etcd_config.InitEtcdCfg(cfg)
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := ecf.Load(ctx); err != nil {
		panic(err)
	}
	go ecf.Watch()

	return cfg
}

func NewNetCfg() *entity.CurrencyNetList {
	cfg := &entity.CurrencyNetList{}
	ecf, err := etcd_config.InitEtcdCfg(cfg)
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := ecf.Load(ctx); err != nil {
		panic(err)
	}
	go ecf.Watch()

	return cfg
}
