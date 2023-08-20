package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"go.uber.org/zap"
)

func GenerateRandomBytes() string {
	keyLength := 256
	key := generateRandomBytes(keyLength / 8)
	keyHex := hex.EncodeToString(key)
	return keyHex
}

func generateRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Failed to generate random bytes")
	}
	return randomBytes
}

func Encrypt(hexKey string, data []byte) []byte {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		Logger.Error("AES密钥处理出错:", zap.Error(err))
		return nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		Logger.Error("创建AES加密器失败", zap.Error(err))
		return nil
	}
	// 使用 AES-CFB 模式创建加密流
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCFBEncrypter(block, iv)
	// 加密数据
	ciphertext := make([]byte, len(data))
	stream.XORKeyStream(ciphertext, data)
	return ciphertext
}

func Decrypt(hexKey string, data []byte) []byte {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		Logger.Error("AES密钥处理出错:", zap.Error(err))
		return nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		Logger.Error("创建AES加密器失败", zap.Error(err))
		return nil
	}
	// 使用 AES-CFB 模式创建加密流
	iv := make([]byte, aes.BlockSize)
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	decrypter.XORKeyStream(decrypted, data)
	return decrypted
}
