package core

import (
	"bufio"
	"github.com/kN6jq/fingerScan/pkg/logger"
	"os"
	"strings"
)

// LoadURLsFromFile 从文件加载URL列表
func LoadURLsFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("读取文件失败: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if url = strings.TrimSpace(url); url != "" {
			if !strings.Contains(url, "http") {
				url = "https://" + url
			}
			urls = append(urls, url)
		}
	}

	return urls
}
