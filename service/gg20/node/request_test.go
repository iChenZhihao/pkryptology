package node

import (
	"encoding/json"
	"fmt"
	"github.com/coinbase/kryptology/service/respvo"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	requestUrl := "http://localhost:8080"
	params := url.Values{}
	str := SecretStrToBase64Str("hello")
	params.Add(DkgSecretKey, str)
	sprintf := fmt.Sprintf("%s?%s", requestUrl, params.Encode())
	fmt.Println(sprintf)
}

func TestReadResp(t *testing.T) {
	var respBytes = []byte{123, 34, 99, 111, 100, 101, 34, 58, 50, 48, 48, 44, 34, 115, 117, 99, 99, 101, 115, 115, 34, 58, 116, 114, 117, 101, 44, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 34, 44, 34, 100, 97, 116, 97, 34, 58, 110, 117, 108, 108, 125}

	var response respvo.Response
	_ = json.Unmarshal(respBytes, &response)
	fmt.Printf("%v\n", response)
}
