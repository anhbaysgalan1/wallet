package wallet_handles

import (
	"github.com/gin-gonic/gin"
	"tp_wallet/internal/wallet/dto"
	"tp_wallet/pkg/hcode"
)

func (h *Handles) BalanceGet(g *gin.Context) {
	var (
		err  error
		req  dto.UidReq
		resp *dto.AccountGetResp
	)
	//uid := h.GetUid(g)
	//if uid == 0 {
	//	return
	//}
	//req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	resp, err = h.walletSrv.BalanceGet(g, &req)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseData(g, resp)
	return
}

func (h *Handles) H2OCash(g *gin.Context) {
	var (
		err error
		req dto.TransferCashReq
	)
	//uid := h.GetUid(g)
	//if uid == 0 {
	//	return
	//}
	//req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	_, err = h.walletSrv.TransferH2OCash(g, &req)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseSuccess(g)
	return
}

func (h *Handles) NftCash(g *gin.Context) {
	var (
		err error
		req dto.NftCashReq
	)
	//uid := h.GetUid(g)
	//if uid == 0 {
	//	return
	//}
	//req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	_, err = h.walletSrv.NftCash(g, &req)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseSuccess(g)
	return
}

func (h *Handles) GetSysTransferAddr(g *gin.Context) {
	var (
		err    error
		result *dto.AddrResp
	)
	result, err = h.walletSrv.GetSysTransferAddr(g, nil)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseData(g, result)
	return
}

func (h *Handles) NftGetByUid(g *gin.Context) {
	var (
		err  error
		req  dto.NftGetByUidReq
		resp *dto.NftInfoS
	)
	//uid := h.GetUid(g)
	//if uid == 0 {
	//	return
	//}
	//req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	resp, err = h.walletSrv.NftGetByUid(g, &req)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseData(g, resp)
	return
}

func (h *Handles) GetSysContractAddr(g *gin.Context) {
	var (
		err  error
		req  dto.ContractType
		resp *dto.AddrResp
	)
	//uid := h.GetUid(g)
	//if uid == 0 {
	//	return
	//}
	//req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	resp, err = h.walletSrv.GetSysContractAddr(g, &req)
	if err != nil {
		h.ResponseErr(g, err)
		return
	}
	h.ResponseData(g, resp)
	return
}
