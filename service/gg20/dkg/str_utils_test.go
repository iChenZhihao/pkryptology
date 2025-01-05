package dkg

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	secret := "hello"
	bytes := SecretStrToBytes(secret)
	base64Str := SecretStrToBase64Str(secret)
	decodeSecret := Base64DecodeSecret(base64Str)
	fmt.Printf("1:%v\n2:%v\n3:%v\n", bytes, base64Str, decodeSecret)
}
