package adapter

import (
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/adapter/grpc"
	"tp_wallet/internal/wallet/adapter/http"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

type Server struct {
	conf      *common.ConfigWalletBusiness
	wallet    dto.WalletSrvServer
	walletJob dto.WalletJob
}

func NewSrv(ws dto.WalletSrvServer, job dto.WalletJob) *Server {
	s := &Server{
		wallet:    ws,
		walletJob: job,
	}
	s.conf = config.WalletBusiness
	s.Init()
	return s
}

func (s *Server) Init() {
	// 服务初始化任务放这里，比如定时任务的开启
	config.NewWalletBusiness()
	config.NewBlockBusiness()
	config.NewConfigFee()
	config.NewConfigCurrency()
	s.walletJob.Run()
}

func (s *Server) RunApp() {
	// 这里必须区块链先注册，然后才能监听
	//s.walletJob.Run()
	http.NewHttp(s.wallet)
	grpc.NewGrpc(s.wallet)
	tool.QuitSignal(func() {
		log.GetLogger().Info("tp_wallet server http exit success", zap.Any("addr", config.WalletBusiness.HttpAddr))
		log.GetLogger().Info("tp_wallet server grpc exit success", zap.Any("addr", config.WalletBusiness.GrpcAddr))
	})
}

func (s *Server) Close() {
	// 处理服务关闭任务
	s.walletJob.Close()
}
