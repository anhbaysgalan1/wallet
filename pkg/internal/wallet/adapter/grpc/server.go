package grpc

import (
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"tp_wallet/config"
	"tp_wallet/pkg/log"
)

func NewGrpc(ws walletPb.WalletSrvServer) {
	grpcServer := grpc.NewServer()
	walletPb.RegisterWalletSrvServer(grpcServer, ws)
	lis, err := net.Listen("tcp", config.WalletBusiness.GrpcAddr)
	if err != nil {
		panic(err)
	}
	log.GetLogger().Info("tp_wallet grpc start success", zap.String("addr:", config.WalletBusiness.GrpcAddr))
	go grpcServer.Serve(lis)
}
