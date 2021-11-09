package hcode

const (
	EN = "en"
	ZH = "zh"
)

// LanguageLen 要支持多少种语言
const LanguageLen = 2

var (
	OK                   = addCode(200)  // 正确
	ErrServer            = addCode(500)  // 服务器内部错误
	ErrReqLimit          = addCode(-501) // 限速
	ErrParameter         = addCode(-502) // 参数错误
	ErrInternalParameter = addCode(-503) // 内部参数错误
	ErrInternalDb        = addCode(-504) // 内部参数错误DB
	ErrInternalCache     = addCode(-505) // 内部参数错误Redis

	ErrAddressGet = addCode(-550) // 地址获取失败

)

func init() {
	addDescription(OK, map[string]string{
		EN: "Success",
		ZH: "成功"})
	addDescription(ErrServer, map[string]string{
		EN: "Server internal error",
		ZH: "服务器内部错误"})
	addDescription(ErrReqLimit, map[string]string{
		EN: "Requests are too frequent",
		ZH: "请求过于频繁，请稍后操作"})
	addDescription(ErrParameter, map[string]string{
		EN: "Parameter error",
		ZH: "参数错误"})
	addDescription(ErrInternalParameter, map[string]string{
		EN: "Internal parameter error",
		ZH: "内部参数错误"})
	addDescription(ErrInternalDb, map[string]string{
		EN: "Internal parameter error",
		ZH: "数据库执行错误"})
	addDescription(ErrInternalCache, map[string]string{
		EN: "Internal parameter error",
		ZH: "缓存执行错误"})

	addDescription(ErrAddressGet, map[string]string{
		EN: "Failed to get address",
		ZH: "地址获取失败"})
}
