package tool

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	var data = make(map[string]interface{})
	data["uid"] = 1183782
	data["cid"] = 2

	//data["uid"] = 254911
	//data["cid"] = 2
	//data["power"] = 1000
	//data["remark"] = "添加用户备注1"

	//data["uid"] = 2549
	//data["cid"] = 2

	data["timestamp"] = fmt.Sprint(GetTimeUnixMilli())
	signStart := JoinStringsInASCII(data, "&", false, "sign")
	signEnd := fmt.Sprintf("%s&sign=%s", signStart, "OrlBruFyyBNzucOy")
	md5Key := GetMD5Encode(signEnd)
	data["sign"] = md5Key
	s, _ := json.Marshal(data)
	fmt.Println(string(s))
	fmt.Println(Sign(data, "OrlBruFyyBNzucOy"))
}
