package http

import (
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"tp_wallet/config"
	"tp_wallet/internal/wallet/adapter/http/router"
	"tp_wallet/internal/wallet/adapter/http/wallet_handles"
	"tp_wallet/pkg/log"
)

func NewHttp(wSrv walletPb.WalletSrvServer) {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	wh := wallet_handles.NewHandles(wSrv)
	router.SetRouters(g, wh)
	server := &http.Server{
		Addr:           config.WalletBusiness.HttpAddr,
		Handler:        g,
		ReadTimeout:    time.Duration(config.WalletBusiness.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WalletBusiness.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.GetLogger().Info("tp_wallet server start success", zap.Any("addr:", config.WalletBusiness.HttpAddr))
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}
