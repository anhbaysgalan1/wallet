package adapter

import (
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.uber.org/zap"
	"tp_wallet/config"
	"tp_wallet/internal/common"
	"tp_wallet/internal/wallet/adapter/grpc"
	"tp_wallet/internal/wallet/adapter/http"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/pkg/database/redis"
	"tp_wallet/pkg/log"
	"tp_wallet/pkg/tool"
)

type Server struct {
	conf      *common.ConfigWalletBusiness
	wallet    walletPb.WalletSrvServer
	walletJob dto.WalletJob
}

func NewSrv(rds *redis.Client, ws walletPb.WalletSrvServer, job dto.WalletJob) *Server {
	s := &Server{
		wallet:    ws,
		walletJob: job,
	}
	s.conf = config.WalletBusiness
	s.Init(rds)
	return s
}

func (s *Server) Init(rds *redis.Client) {
	// 服务初始化任务放这里，比如定时任务的开启
	config.NewWalletBusiness()
	config.NewBlockBusiness()
	config.NewConfigFee()
	config.NewConfigCurrency()
	// 初始化雪花算法
	if snowFlakeId, err := tool.NewWorkerIdTool(rds); err != nil {
		log.GetLogger().Error("[Server Init] tool.NewWorkerIdTool failed", zap.Error(err))
		panic(err)
	} else {
		if err := tool.NewWorker(snowFlakeId); err != nil {
			log.GetLogger().Error("[Server Init] tool.NewWorker failed", zap.Int64("id", snowFlakeId),
				zap.Error(err))
			panic(err)
		}
		log.GetLogger().Info("tool.NewWorkerIdTool snow flake init success", zap.Int64("id", snowFlakeId))
	}
	// job服务暂时关闭
	s.walletJob.Run()
}

func (s *Server) RunApp(cancel func()) {
	// 这里必须区块链先注册，然后才能监听
	//s.walletJob.Run()
	http.NewHttp(s.wallet)
	grpc.NewGrpc(s.wallet)
	tool.QuitSignal(func() {
		cancel()
		s.Close()
		log.GetLogger().Info("tp_wallet server http exit success", zap.Any("addr", config.WalletBusiness.HttpAddr))
		log.GetLogger().Info("tp_wallet server grpc exit success", zap.Any("addr", config.WalletBusiness.GrpcAddr))
	})
}

func (s *Server) Close() {
	// 处理服务关闭任务
	s.walletJob.Close()
}
