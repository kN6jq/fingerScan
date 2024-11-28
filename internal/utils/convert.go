package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// GetCurrentPath 获取当前程序路径
func GetCurrentPath() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return path
}

// RemoveDuplicates 去除字符串切片中的重复项和空值
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if entry != "" {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
	}
	return list
}

// ToJSON 将对象转换为JSON字符串
func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// FromJSON 从JSON字符串解析对象
func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// ParseInt 安全的字符串转整数
func ParseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// JoinPath 连接路径
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// EnsureDir 确保目录存在
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
