package app

import (
	"github.com/gin-gonic/gin"
	"github.com/leaf-rain/wallet/pkg/convert"
	"github.com/leaf-rain/wallet/pkg/hcode"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
	"net/http"
)

func GetUid(ctx *gin.Context) int {
	g, exist := ctx.Get("uid")
	if !exist {
		return 0
	} else {
		return g.(int)
	}
}

func GetCid(ctx *gin.Context) int {
	g, exist := ctx.Get("cid")
	if !exist {
		return 0
	} else {
		return g.(int)
	}
}
func GetAttrId(g *gin.Context) int {
	key, ok := g.Get("aid")
	if !ok {
		return 0
	}
	return key.(int)
}

func ResponseErr(g *gin.Context, err error) {
	code := hcode.Cause(err)
	data := gin.H{
		"code": code.Code(),
		"data": "",
		"msg":  code.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseErr", zap.Any("res", data))
	g.JSON(200, data)
}
func ResponseSuccess(g *gin.Context) {
	info := gin.H{
		"code": hcode.OK,
		"data": "",
		"msg":  hcode.OK.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseData", zap.Any("res", info))
	g.JSON(200, info)
}

func ResponseData(g *gin.Context, data interface{}) {
	info := gin.H{
		"code": hcode.OK,
		"data": data,
		"msg":  hcode.OK.Message(g.GetHeader("lang")),
	}
	log.GetLogger().Debug("ResponseData", zap.Any("res", info))
	g.JSON(200, info)
}

type Pager struct {
	// 页码
	Page int `json:"page"`
	// 每页数量
	PageSize int `json:"page_size"`
	// 总行数
	TotalRows int `json:"total_rows"`
}

func ToResponseList(g *gin.Context, list interface{}, totalRows int) {
	g.JSON(http.StatusOK, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(g),
			PageSize:  GetPageSize(g),
			TotalRows: totalRows,
		},
	})
}

func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("page_size")).MustInt()
	//if pageSize <= 0 {
	//	return global.AppSetting.DefaultPageSize
	//}
	//if pageSize > global.AppSetting.MaxPageSize {
	//	return global.AppSetting.MaxPageSize
	//}
	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	var result = 0
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}

type Pagination struct {
	Limit     int64 `json:"limit"`
	Offset    int64 `json:"offset"`
	TotalRows int64 `json:"total_rows"`
}

func GetPagination(c *gin.Context) *Pagination {
	var page = Pagination{}
	page.Limit = convert.StrTo(c.Query("limit")).MustInt64()
	page.Offset = convert.StrTo(c.Query("offset")).MustInt64()
	if page.Limit == 0 {
		return nil
	}
	return &page
}
