package node

import (
	"encoding/json"
	"fmt"
	"github.com/coinbase/kryptology/service/respvo"
	"github.com/coinbase/kryptology/service/utils"
	"github.com/golang/glog"
	"time"
)

var httpClient = utils.NewHTTPClient(20 * time.Second)

func DoSendStartDkg(requestUrl string) error {
	responseByte, err := httpClient.Get(requestUrl, nil, nil)
	var response respvo.Response
	err = json.Unmarshal(responseByte, &response)
	if err != nil {
		glog.Errorf("unmarshal error: %v", err.Error())
		return err
	}
	if !response.Success {
		glog.Errorf("not success: %v", err.Error())
		return fmt.Errorf("DoSendStartDkg Fail")
	}
	return nil
}

// DoSendBroadcastRound1 sends a message to a single node via HTTP POST.
func DoSendBroadcastRound1(url string, message DkgRound1Recv) error {
	responseByte, err := httpClient.Post(url, nil, message)
	//glog.Infof("responseByte: %v", responseByte)
	if err != nil {
		return err
	}
	var response respvo.Response
	err = json.Unmarshal(responseByte, &response)
	if err != nil {
		glog.Errorf("unmarshal error: %v", err.Error())
		return err
	}
	if !response.Success {
		glog.Errorf("not success: %v", err.Error())
		return fmt.Errorf("DoSendStartDkg Fail")
	}
	return nil
}

// DoSendBroadcastRound2 sends a message to a single node via HTTP POST.
func DoSendBroadcastRound2(url string, message DkgRound2Recv) error {
	responseByte, err := httpClient.Post(url, nil, message)
	if err != nil {
		return err
	}
	var response respvo.Response
	err = json.Unmarshal(responseByte, &response)
	if err != nil {
		glog.Errorf("unmarshal error: %v", err.Error())
		return err
	}
	if !response.Success {
		glog.Errorf("not success: %v", err.Error())
		return fmt.Errorf("DoSendStartDkg Fail")
	}
	return nil
}

// DoSendBroadcastRound3 sends a message to a single node via HTTP POST.
func DoSendBroadcastRound3(url string, message DkgRound3Recv) error {
	responseByte, err := httpClient.Post(url, nil, message)
	if err != nil {
		return err
	}
	var response respvo.Response
	err = json.Unmarshal(responseByte, &response)
	if err != nil {
		glog.Errorf("unmarshal error: %v", err.Error())
		return err
	}
	if !response.Success {
		glog.Errorf("not success: %v", err.Error())
		return fmt.Errorf("DoSendStartDkg Fail")
	}
	return nil
}
