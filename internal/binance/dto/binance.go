package dto

import (
	"github.com/leaf-rain/wallet/internal/binance/consts"
	tool2 "github.com/leaf-rain/wallet/pkg/tool"
	"net/url"
)

type Parameters struct {
	RecvWindow int64 `json:"recvWindow"` // 时间差
	Timestamp  int64 `json:"timestamp"`  // 时间戳
}

func (a *Parameters) ToString() string {
	if a.RecvWindow <= 0 {
		a.RecvWindow = consts.RecvWindow
	}
	if a.Timestamp <= 0 {
		a.Timestamp = tool2.GetTimeUnixMilli()
	}
	var values = url.Values{}
	values.Add("recvWindow", tool2.Int64ToString(a.RecvWindow))
	values.Add("timestamp", tool2.Int64ToStr(a.Timestamp))
	return values.Encode()
}
