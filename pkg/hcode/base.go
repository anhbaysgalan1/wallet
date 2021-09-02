package hcode

const (
	EN = "en"
	ZH = "zh"
)

// LanguageLen 要支持多少种语言
const LanguageLen = 2

var (
	OK        = addCode(200) // 正确
	ReqLimit  = addCode(201) // 限速
	ServerErr = addCode(300) // 服务器内部错误
)

func init() {
	addDescription(OK, map[string]string{
		EN: "Success",
		ZH: "成功"})
	addDescription(ReqLimit, map[string]string{
		EN: "Requests are too frequent",
		ZH: "请求过于频繁，请稍后操作"})
	addDescription(ServerErr, map[string]string{
		EN: "Server internal error",
		ZH: "服务器内部错误"})
}
