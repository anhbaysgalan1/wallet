package tool

import (
	"fmt"
	"gluttonous/consts"
	"strconv"
)

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func Uint32ToString(num uint32) string {
	return strconv.FormatUint(uint64(num), 10)
}

func FloatPriceToInt64(price float64) int64 {
	return int64(price * consts.CoinConversion)
}

func Int64PriceToFloat(price int64) float64 {
	return float64(price) / consts.CoinConversion
}

//取6位精度并转换成int64
func To6Decimal(value float64) int64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", value), 64)
	return int64(value * consts.CoinConversion)
}

//取2位精度并转换成int64
func To2Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
