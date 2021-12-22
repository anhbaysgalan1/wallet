package wallet_handles

import (
	hcode "git.huoys.com/chain-game/rowing_proto/common/error"
	walletPb "git.huoys.com/chain-game/rowing_proto/wallet"
	"github.com/gin-gonic/gin"
)

func (h *Handles) BalanceGet(g *gin.Context) {
	var (
		err  error
		req  walletPb.UidReq
		resp *walletPb.AccountGetResp
	)
	uid := h.GetUid(g)
	if uid == 0 {
		return
	}
	req.Uid = uid
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

func (h *Handles) CurrencyCash(g *gin.Context) {
	var (
		err error
		req walletPb.TransferCashReq
	)
	uid := h.GetUid(g)
	if uid == 0 {
		return
	}
	req.Uid = uid
	err = g.BindJSON(&req)
	if err != nil {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	_, err = h.walletSrv.TransferCurrencyCash(g, &req)
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
		req walletPb.NftCashReq
	)
	uid := h.GetUid(g)
	if uid == 0 {
		return
	}
	req.Uid = uid
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
		req    *walletPb.CurrencyReq
		result *walletPb.AddrResp
	)
	err = g.BindJSON(&req)
	if err != nil || req.Currency == "" {
		h.ResponseErr(g, hcode.ErrParameter)
		return
	}
	result, err = h.walletSrv.GetSysTransferAddr(g, req)
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
		req  walletPb.NftGetByUidReq
		resp *walletPb.NftInfoS
	)
	uid := h.GetUid(g)
	if uid == 0 {
		return
	}
	req.Uid = uid
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
		req  walletPb.ContractType
		resp *walletPb.AddrResp
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
