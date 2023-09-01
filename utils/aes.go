package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func GenerateRandomBytes() string {
	key := generateRandomBytes(32)
	fmt.Printf("%x\n", key)
	hexKey := hex.EncodeToString(key)
	return hexKey
}

func generateRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Failed to generate random bytes")
	}
	return randomBytes
}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 移除
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncrypt 加密
func Encrypt(base64key string, data []byte) []byte {
	key, err := hex.DecodeString(base64key)
	fmt.Println(len(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	fmt.Println(blockSize)
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	fmt.Println(base64.StdEncoding.EncodeToString(key[:blockSize]))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted
}

// AesDecrypt 解密
func Decrypt(base64Key string, data []byte) []byte {
	key, err := hex.DecodeString(base64Key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil
	}
	return crypted
}

//func Encrypt(hexKey string, data []byte) []byte {
//	key, err := base64.StdEncoding.DecodeString(hexKey)
//	if err != nil {
//		Logger.Error("AES密钥处理出错:", zap.Error(err))
//		return nil
//	}
//	// create cipher
//	c, err := aes.NewCipher(key)
//	if err != nil {
//		Logger.Error("创建AES Cipher失败", zap.Error(err))
//		return nil
//	}
//	// allocate space for ciphered data
//	out := make([]byte, len(data))
//
//	// encrypt
//	c.Encrypt(out, data)
//	// return hex string
//	return out
//}
//
//func Decrypt(hexKey string, data []byte) []byte {
//	key, err := base64.StdEncoding.DecodeString(hexKey)
//	if err != nil {
//		Logger.Error("AES密钥处理出错:", zap.Error(err))
//		return nil
//	}
//	c, err := aes.NewCipher(key)
//	if err != nil {
//		Logger.Error("创建AES Cipher失败", zap.Error(err))
//		return nil
//	}
//	out := make([]byte, len(data))
//	c.Decrypt(out, data)
//	return out
//}
