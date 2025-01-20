package dkg

import (
	"encoding/base64"
	"fmt"
)

// SecretStrToBase64Str 对secret串进行base64编码
func SecretStrToBase64Str(secret string) string {
	if secret == "" {
		return ""
	}
	data := []byte(secret)
	return base64.StdEncoding.EncodeToString(data)
}

func SecretStrToBytes(secret string) []byte {
	if secret == "" {
		return nil
	}
	return []byte(secret)
}

func OtherNodeStartDkgUrl(nodeAddr string) string {
	return fmt.Sprintf("http://%s/dkg/round1", nodeAddr)
}

func GetDkgRound1BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("http://%s/dkg/round1/recv", nodeAddr)
}

func GetDkgRound2BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("http://%s/dkg/round2/recv", nodeAddr)
}

func GetDkgRound3BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("http://%s/dkg/round3/recv", nodeAddr)
}

func Base64DecodeSecret(secretStr string) []byte {
	if len(secretStr) == 0 {
		return nil
	}
	decodedData, _ := base64.StdEncoding.DecodeString(secretStr)
	return decodedData
}
