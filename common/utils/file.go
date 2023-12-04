package utils

import (
	"go.uber.org/zap"
	"orange-clipboard/common/resource"
	"os"
)

func ReadFile(fileName string) []byte {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		resource.Logger.Error("读取文件失败", zap.String("filename", fileName), zap.Error(err))
		return nil
	}
	return bytes
}
