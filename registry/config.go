package registry

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/youpenglai/goutils/httptool"
)

var (
	ErrConfKeyNotExists = errors.New("config key not exists")
)

type ConsulKvMetaDataResponseItem struct {
	CreateIndex int64  `json:"CreateIndex"`
	ModifyIndex int64  `json:"ModifyIndex"`
	LockIndex   int64  `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	Session     string `json:"Session"`
}

func ReadConf(confKey string) (data []byte, err error) {
	queryUrl := fmt.Sprintf("http://%s:%d/v1/kv/%s", consulIp, consulPort, confKey)
	var kvMetaData []byte
	_, kvMetaData, err = httptool.HttpGet(queryUrl, nil)
	if err != nil {
		return
	}

	var kvMeta []*ConsulKvMetaDataResponseItem
	if err = json.Unmarshal(kvMetaData, &kvMeta); err != nil {
		return
	}

	data, err = base64.StdEncoding.DecodeString(kvMeta[0].Value)

	return
}

// TODO: 增加其他内容
