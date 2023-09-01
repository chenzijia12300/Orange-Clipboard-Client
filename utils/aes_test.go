package utils

import (
	"encoding/base64"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	t.Log(GenerateRandomBytes())
}

func TestEncrypt(t *testing.T) {
	base64Key := "1b18a76c028b970b5f27e6eed5eeacd0cd9f26e34c3a6a319c2a6c1761caed7f"
	data := Encrypt(base64Key, []byte("go"))
	t.Log(string(Decrypt(base64Key, data)))
	t.Log(base64.StdEncoding.EncodeToString(data))
}
