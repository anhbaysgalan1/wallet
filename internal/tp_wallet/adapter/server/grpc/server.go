package grpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"tp_wallet/config"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/pkg/log"
)

func NewGrpc(ws dto.WalletSrvServer) {
	grpcServer := grpc.NewServer()
	dto.RegisterWalletSrvServer(grpcServer, ws)
	lis, err := net.Listen("tcp", config.WalletBusiness.GrpcAddr)
	if err != nil {
		panic(err)
	}
	log.GetLogger().Info("tp_wallet grpc start success", zap.String("addr:", config.WalletBusiness.GrpcAddr))
	go grpcServer.Serve(lis)
}
