package entity

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/leaf-rain/wallet/internal/currency_manager/dto"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
)

type Currency struct {
	Name   string `json:"name" bson:"name" yaml:"name"`       // 币种名称
	Status int64  `json:"status" bson:"status" yaml:"status"` // 二进制开关，0 提币，1 充值， 2 划转
}

func (c *Currency) ToPb() *dto.Currency {
	return &dto.Currency{
		Name:   c.Name,
		Status: c.Status,
	}
}

type CurrencyList []Currency

func (a *CurrencyList) GetConfigKey() string {
	return "/github.com/leaf-rain/wallet/Currency"
}

func (a *CurrencyList) LoadConfig(value []byte) {
	if err := yaml.Unmarshal(value, &a); err != nil {
		log.GetLogger().Error("Unmarshal Yaml Failed", zap.Any("value", value))
		return
	}
	fmt.Println("LoadConfig", "value: ", string(value))
	fmt.Println("LoadConfig", "config: ", *a)
}

type CurrencyNet struct {
	Name      string `json:"name,omitempty" bson:"name" yaml:"name"`                   // 网络名称
	SpareName string `json:"spare_name,omitempty" bson:"spare_name" yaml:"spare_name"` // 备用名称（目前是传给区块链）
	MasterCy  string `json:"master_cy,omitempty" bson:"master_cy" yaml:"master_cy"`    // 转账币种
	AddressCy string `json:"address_cy" bson:"address_cy" yaml:"address_cy"`           // 获取地址币种
	Status    bool   `json:"status,omitempty" bson:"status" yaml:"status"`             // 是否启用
	AdviseGas string `json:"advise_gas" bson:"advise_gas" yaml:"advise_gas"`           // 建议gas
	MinGas    string `json:"min_gas" bson:"min_gas" yaml:"min_gas"`                    // 最小手续费
	MaxGas    string `json:"max_gas" bson:"max_gas" yaml:"max_gas"`                    // 最大手续费
	Weight    int32  `json:"weight,omitempty" bson:"weight" yaml:"weight"`             // 排序权重
}

func (c *CurrencyNet) ToPb() *dto.Net {
	return &dto.Net{
		Name:      c.Name,
		SpareName: c.SpareName,
		MasterCy:  c.MasterCy,
		AddressCy: c.AddressCy,
		Status:    c.Status,
		AdviseGas: c.AdviseGas,
		MinGas:    c.MinGas,
		MaxGas:    c.MaxGas,
		Weight:    c.Weight,
	}
}

type CurrencyNetList []CurrencyNet

func (a *CurrencyNetList) GetConfigKey() string {
	return "/github.com/leaf-rain/wallet/Currency_net"
}

func (a *CurrencyNetList) LoadConfig(value []byte) {
	if err := yaml.Unmarshal(value, &a); err != nil {
		log.GetLogger().Error("Unmarshal Yaml Failed", zap.Any("value", value))
		return
	}
	fmt.Println("LoadConfig", "value: ", string(value))
	fmt.Println("LoadConfig", "config: ", *a)
}
