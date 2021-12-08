package wallet_handles

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/pkg/hcode"
	"tp_wallet/pkg/log"
)

type Handles struct {
	walletSrv dto.WalletSrvServer
}

func NewHandles(wSrv dto.WalletSrvServer) *Handles {
	return &Handles{walletSrv: wSrv}
}

func (h *Handles) GetUid(g *gin.Context) uint64 {
	key, ok := g.Get("uid")
	if !ok {
		h.ResponseErr(g, hcode.ErrParameter)
		return 0
	}
	return key.(uint64)
}

func (h *Handles) ResponseErr(g *gin.Context, err error) {
	code := hcode.Cause(err)
	data := gin.H{
		"code": code.Code(),
		"data": "",
		"msg":  code.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseErr", zap.Any("res", data))
	g.JSON(200, data)
}

func (h *Handles) ResponseSuccess(g *gin.Context) {
	info := gin.H{
		"code": hcode.OK,
		"data": "",
		"msg":  hcode.OK.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseData", zap.Any("res", info))
	g.JSON(200, info)
}

func (h *Handles) ResponseData(g *gin.Context, data interface{}) {
	info := gin.H{
		"code": hcode.OK,
		"data": data,
		"msg":  hcode.OK.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseData", zap.Any("res", info))
	g.JSON(200, info)
}
