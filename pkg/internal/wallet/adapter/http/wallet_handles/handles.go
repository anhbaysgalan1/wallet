package wallet_handles

import (
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	kgin "git.huoys.com/middle-end/kratos/pkg/net/http/gin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"tp_wallet/pkg/log"
)

type Handles struct {
	walletSrv walletPb.WalletSrvServer
}

func NewHandles(wSrv walletPb.WalletSrvServer) *Handles {
	return &Handles{walletSrv: wSrv}
}

func (h *Handles) GetUid(g *gin.Context) uint64 {
	key, ok := g.Request.URL.Query()["userid"]
	if !ok || len(key) != 1 {
		h.ResponseErr(g, hcode.ErrUserNotFound)
		return 0
	}
	uid, err := strconv.ParseUint(key[0], 10, 64)
	if uid == 0 {
		log.GetLogger().Error("[GetUid] failed", zap.Error(err))
		h.ResponseErr(g, hcode.ErrUserNotFound)
		return 0
	}
	return uid

}

func (h *Handles) ResponseErr(g *gin.Context, err error) {
	resp := kgin.TOJSON(nil, err)
	g.JSON(http.StatusOK, resp)
}

func (h *Handles) ResponseSuccess(g *gin.Context) {
	resp := kgin.TOJSON(nil, nil)
	g.JSON(http.StatusOK, resp)
}

func (h *Handles) ResponseData(g *gin.Context, data interface{}) {
	resp := kgin.TOJSON(data, nil)
	g.JSON(http.StatusOK, resp)
}
