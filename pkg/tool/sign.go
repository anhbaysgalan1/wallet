package tool

import (
	"fmt"
	"gluttonous/pkg/hcode"
	"sort"
	"strings"
)

func Sign(data map[string]interface{}, key string) error {
	reqSign, ok := data["sign"]
	if !ok {
		return hcode.ParameterErr
	}
	signStart := JoinStringsInASCII(data, "&", false, "sign")
	signEnd := fmt.Sprintf("%s&sign=%s", signStart, key)
	endSign := GetMD5Encode(signEnd)
	if reqSign != endSign {
		return hcode.SignErr
	}
	return nil
}

//JoinStringsInASCII 按照规则，参数名ASCII码从小到大排序后拼接
//data 待拼接的数据
//sep 连接符
//includeEmpty 是否包含空值，true则包含空值，否则不包含，注意此参数不影响参数名的存在
//exceptKeys 被排除的参数名，不参与排序及拼接
func JoinStringsInASCII(data map[string]interface{}, sep string, includeEmpty bool, exceptKeys ...string) string {
	var list []string
	m := make(map[string]int)
	if len(exceptKeys) > 0 {
		for _, except := range exceptKeys {
			m[except] = 1
		}
	}
	for k := range data {
		if _, ok := m[k]; ok {
			continue
		}
		value := data[k]
		if !includeEmpty && value == "" {
			continue
		}
		list = append(list, fmt.Sprintf("%s=%v", k, value))
	}
	sort.Strings(list)
	return strings.Join(list, sep)
}
