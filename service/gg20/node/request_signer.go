package node

import (
	"encoding/json"
	"fmt"
	"github.com/coinbase/kryptology/service/respvo"
	"github.com/coinbase/kryptology/service/utils"
	"github.com/golang/glog"
	"time"
)

var signerHttpClient = utils.NewHTTPClient(15 * time.Second)

func DoSendAskForCosignerCandidate(nodeAddress, workId string) (*CandidateInfo, error) {
	post, err := signerHttpClient.Post(nodeAddress, nil, workId)
	if err != nil {
		return nil, err
	}
	var response respvo.Response
	err = json.Unmarshal(post, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("返回数据类型不匹配")
	}
	return &CandidateInfo{Id: uint32(info["id"].(float64)), WorkId: info["workId"].(string)}, nil
}

func DoPostWithoutRespData(requestUrl string, data interface{}) error {
	post, err := httpClient.Post(requestUrl, nil, data)
	if err != nil {
		return err
	}
	var resp respvo.Response
	err = json.Unmarshal(post, &resp)
	if err != nil {
		glog.Errorf("unmarshal error: %v", err.Error())
		return err
	}
	if !resp.Success {
		glog.Errorf("not success: %v", err.Error())
		return fmt.Errorf("DoSendStartDkg Fail")
	}
	return nil
}
