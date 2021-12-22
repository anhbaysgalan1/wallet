package router

import (
	"github.com/gin-gonic/gin"
	"tp_wallet/internal/wallet/adapter/http/wallet_handles"
)

func SetRouters(r *gin.Engine, walletHandles *wallet_handles.Handles) {
	SetCorsRouters(r)
	wallet := r.Group("/wallet")
	wallet.POST("/currency/get_balance", walletHandles.BalanceGet)
	wallet.POST("/currency/cash", walletHandles.CurrencyCash)
	wallet.POST("/currency/get_sys_addr", walletHandles.GetSysTransferAddr)
	wallet.POST("/nft/cash", walletHandles.NftCash)
	wallet.POST("/nft/get_sys_addr", walletHandles.GetSysContractAddr)
	wallet.POST("/nft/get_by_uid", walletHandles.NftGetByUid)
}
