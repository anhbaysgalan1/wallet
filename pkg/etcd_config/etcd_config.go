package etcd_config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Config interface {
	GetConfigKey() string
	LoadConfig(value []byte)
}

//todo close
type EtcdCfg struct {
	CFG     Config
	ETCDCli *clientv3.Client
	Watcher clientv3.Watcher
}

func InitEtcdCfg(config Config) (*EtcdCfg, error) {
	ECF := &EtcdCfg{}

	//ecf.CFG = config
	if config == nil {
		return nil, errors.New("Fuck! Your config is nil")
	}
	ECF.CFG = config
	//(*ecf.CFG).GetConfigKey()
	ETCD := os.Getenv("ETCD")
	TLSTR := os.Getenv("TLS")
	TLS, err := strconv.ParseBool(TLSTR)
	if err != nil {
		TLS = false
	}

	var CliConf clientv3.Config

	if TLS {
		etcdCert := "/ca/cert.pem"
		etcdCertKey := "/ca/key.pem"
		etcdCa := "/ca/ca.pem"

		cert, err := tls.LoadX509KeyPair(etcdCert, etcdCertKey)
		if err != nil {
			return nil, err
		}

		caData, err := ioutil.ReadFile(etcdCa)
		if err != nil {
			return nil, err
		}

		ca := x509.NewCertPool()
		ca.AppendCertsFromPEM(caData)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      ca,
		}

		CliConf = clientv3.Config{
			Endpoints:   []string{ETCD},
			DialTimeout: 5 * time.Second,
			TLS:         tlsConfig,
		}
	} else {
		CliConf = clientv3.Config{
			Endpoints:   []string{ETCD},
			DialTimeout: 5 * time.Second,
		}
	}

	cli, err := clientv3.New(CliConf)
	if err != nil {
		return nil, errors.New("Fuck! ETCD can't connect")
	}
	ECF.ETCDCli = cli

	watcher := clientv3.NewWatcher(ECF.ETCDCli)

	ECF.Watcher = watcher
	return ECF, nil
}

func (ecf *EtcdCfg) Watch() {
	wch := ecf.Watcher.Watch(context.TODO(), ecf.CFG.GetConfigKey())

watch:
	for {
		select {
		case c := <-wch:
			if c.Err() != nil {
				//todo 打印错误日志
				break watch
			} else {
				for _, event := range c.Events {
					switch event.Type {
					case mvccpb.PUT:
						log.GetLogger().Debug("Watch event PUT", zap.Any("event Value", event.Kv.Value))
						ecf.CFG.LoadConfig(event.Kv.Value)
						//event.Kv.Value
					case mvccpb.DELETE:
						//todo 处理Delete
						log.GetLogger().Debug("Watch event DELETE", zap.Any("event Value", event.Kv.Value))
					}
				}
			}

		}
	}
	//fmt.Println("re watch")
	//ecf.Watch()
}

func (ecf *EtcdCfg) Load(ctx context.Context) error {
	if resp, err := ecf.ETCDCli.Get(ctx, ecf.CFG.GetConfigKey()); err != nil {
		return err
	} else {
		if resp.Count == 1 {
			ecf.CFG.LoadConfig(resp.Kvs[0].Value)
		} else {
			log.GetLogger().Error("load value Count more than 1", zap.Any("Count", resp.Count))
		}
	}
	return nil
}

func (ecf *EtcdCfg) Save(ctx context.Context, value []byte) error {
	if _, err := ecf.ETCDCli.Put(ctx, ecf.CFG.GetConfigKey(), string(value)); err != nil {
		return err
	}
	return nil
}

func (ecf *EtcdCfg) Close() error {
	if err := ecf.Watcher.Close(); err != nil {
		return err
	}
	if err := ecf.ETCDCli.Close(); err != nil {
		return err
	}
	return nil
}
