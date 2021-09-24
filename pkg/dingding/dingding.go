package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leaf-rain/wallet/pkg/log"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DingReq struct {
	Secret               string
	Urls                 string
	AtMobiles            []string
	AtUserIds            []string
	MinBalanceForMining  int64
	MinBalanceForAirdrop int64
}

var dingReq DingReq

func init() {
	dingReq = DingReq{
		Secret:               "SEC12d98e9eb8a48e056cd7183e23f1368bbaf1a9f5214a514b47ce6f1e149485ea",
		Urls:                 "https://oapi.dingtalk.com/robot/send?access_token=e87ba7b0cf184d30c5db9fafab3553c871124196d9330f5edf57abeef7f0602f",
		AtMobiles:            []string{"17674923763"},
		AtUserIds:            []string{"yeyangfengqi"},
		MinBalanceForMining:  10000,
		MinBalanceForAirdrop: 100000,
	}
	if log.GetLogger() == nil {
		_, _ = log.NewLogger(&log.Options{
			AppName: "gluttonous_test",
			Level:   "info",
		})
	}
}

type ResultForDingDing struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func DingdingSend(message string) (bool, error) {
	var reqUrl = dingReq.Urls + sign(dingReq.Secret)
	var request, err = newRequestBody(dingReq.AtMobiles, dingReq.AtUserIds, message, false)
	if err != nil {
		log.GetLogger().Error("[DingdingSend] newRequestBody failed",
			zap.Any("setting", dingReq),
			zap.Any("message", message),
			zap.Error(err))
		return false, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(request))
	if err != nil {
		log.GetLogger().Error("[DingdingSend] NewRequest failed",
			zap.Any("setting", dingReq),
			zap.Any("message", message),
			zap.Error(err))
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		log.GetLogger().Error("[DingdingSend] request Do failed",
			zap.Any("setting", dingReq),
			zap.Any("message", message),
			zap.Error(err))
		return false, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.GetLogger().Error("[DingdingSend] ioutil.ReadAll failed",
			zap.Any("setting", dingReq),
			zap.Any("message", message),
			zap.Error(err))
		return false, err
	}
	var sta ResultForDingDing
	err = json.Unmarshal(body, &sta)
	if err != nil {
		log.GetLogger().Error("[DingdingSend] json.Unmarshal failed",
			zap.Any("setting", dingReq),
			zap.Any("message", message),
			zap.Error(err))
		return false, err
	}
	if sta.Errcode == 0 || sta.Errmsg == "ok" {
		log.GetLogger().Info("[DingdingSend] success",
			zap.Any("message", message))
		return true, nil
	}
	return false, errors.New("Error response:---" + string(body))
}

func newRequestBody(atm, atu []string, data string, isAtAll bool) (string, error) {
	reqBody := struct {
		At struct {
			AtMobiles []string `json:"atMobiles"`
			AtUserIds []string `json:"atUserIds"`
			IsAtAll   bool     `json:"isAtAll"`
		} `json:"at"`
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
		Msgtype string `json:"msgtype"`
	}{}
	reqBody.Text.Content = data
	reqBody.Msgtype = "text"
	reqBody.At.AtMobiles = atm
	reqBody.At.AtUserIds = atu
	reqBody.At.IsAtAll = isAtAll
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	return string(reqData), nil
}

func sign(secret string) string {
	timestamp := fmt.Sprint(time.Now().UnixNano() / 1000000)
	secStr := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(secStr))
	sum := h.Sum(nil)
	encode := base64.StdEncoding.EncodeToString(sum)
	urlEncode := url.QueryEscape(encode)
	return "&timestamp=" + timestamp + "&sign=" + urlEncode
}
