package node

import (
	"encoding/base64"
	"fmt"
)

const Protocol = "http://"

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
	return fmt.Sprintf("%s%s/dkg/round1", Protocol, nodeAddr)
}

func GetDkgRound1BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/dkg/round1/recv", Protocol, nodeAddr)
}

func GetDkgRound2BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/dkg/round2/recv", Protocol, nodeAddr)
}

func GetDkgRound3BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/dkg/round3/recv", Protocol, nodeAddr)
}

func Base64DecodeSecret(secretStr string) []byte {
	if len(secretStr) == 0 {
		return nil
	}
	decodedData, _ := base64.StdEncoding.DecodeString(secretStr)
	return decodedData
}

func GetAskCosignerCandidateUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/candidate", Protocol, nodeAddr)
}

func GetOtherStartSignUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round1", Protocol, nodeAddr)
}

func GetSignRound1BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round1/recv", Protocol, nodeAddr)
}

func GetSignRound2BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round2/recv", Protocol, nodeAddr)
}

func GetSignRound3BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round3/recv", Protocol, nodeAddr)
}

func GetSignRound4BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round4/recv", Protocol, nodeAddr)
}

func GetSignRound5BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round5/recv", Protocol, nodeAddr)
}

func GetSignRound6BcastUrl(nodeAddr string) string {
	return fmt.Sprintf("%s%s/sign/round6/recv", Protocol, nodeAddr)
}
