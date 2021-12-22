package bsc

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
	"tp_wallet/internal/block_chain/chain/scan/common"
)

func requestExec(method, key, param string) ([]byte, error) {
	time.Sleep(time.Second / 4)
	url := common.ChainNetUrl + method + param + "&apikey=" + key
	tl := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tl}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
