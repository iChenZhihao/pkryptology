package dkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	response2 "github.com/coinbase/kryptology/service/response"
	"github.com/golang/glog"
	"io"
	"net/http"
	//"net/url"
	"time"
)

func DoSendStartDkg(requestUrl string) error {
	client := &http.Client{
		Timeout: 50 * time.Second, // 设置超时时间
	}
	//params := url.Values{}
	//params.Add(DkgSecretKey, secretBase64)
	//sprintf := fmt.Sprintf("%s?%s", requestUrl, params.Encode())
	//glog.Infof("final url: %s\n", sprintf)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	responseByte, err := io.ReadAll(resp.Body)
	//glog.Infof("responseByte: %v", responseByte)
	if err != nil {
		return err
	}
	var response response2.Response
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
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	client := &http.Client{
		Timeout: 21 * time.Second, // 设置超时时间
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	responseByte, err := io.ReadAll(resp.Body)
	//glog.Infof("responseByte: %v", responseByte)
	if err != nil {
		return err
	}
	var response response2.Response
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
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	client := &http.Client{
		Timeout: 21 * time.Second, // 设置超时时间
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	responseByte, err := io.ReadAll(resp.Body)
	//glog.Infof("responseByte: %v", responseByte)
	if err != nil {
		return err
	}
	var response response2.Response
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
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	client := &http.Client{
		Timeout: 21 * time.Second, // 设置超时时间
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	responseByte, err := io.ReadAll(resp.Body)
	//glog.Infof("responseByte: %v", responseByte)
	if err != nil {
		return err
	}
	var response response2.Response
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
