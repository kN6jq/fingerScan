// Package utils 提供HTTP相关的工具函数
package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/twmb/murmur3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	faviconTimeout = 8 * time.Second
	base64LineLen  = 76
)

var (
	jsRedirectPatterns = []string{
		`(window|top)\.location\.href = ['"](.*?)['"]`,
		`redirectUrl = ['"](.*?)['"]`,
		`<meta.*?http-equiv=.*?refresh.*?url=(.*?)>`,
	}
)

// RandomUserAgent 返回随机User-Agent
func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		// ... 其他User-Agent
	}
	return userAgents[rand.Intn(len(userAgents))]
}

// CalculateFaviconHash 计算favicon的哈希值
func CalculateFaviconHash(url string) string {
	if url == "" {
		return "0"
	}

	favicon, err := fetchFavicon(url)
	if err != nil {
		return "0"
	}

	encodedFavicon := encodeBase64WithLineBreaks(favicon)
	return calculateMurmurHash(encodedFavicon)
}

// fetchFavicon 获取favicon图标
func fetchFavicon(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: faviconTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("favicon request failed: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// encodeBase64WithLineBreaks 使用换行符对数据进行base64编码
func encodeBase64WithLineBreaks(data []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(data)
	var buffer bytes.Buffer

	for i := 0; i < len(encoded); i++ {
		buffer.WriteByte(encoded[i])
		if (i+1)%base64LineLen == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')

	return buffer.Bytes()
}

// calculateMurmurHash 计算MurmurHash
func calculateMurmurHash(data []byte) string {
	h32 := murmur3.New32()
	if _, err := h32.Write(data); err != nil {
		return "0"
	}
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

// ExtractJSURLs 提取JavaScript重定向URL
func ExtractJSURLs(body, baseURL string) []string {
	var urls []string
	for _, pattern := range jsRedirectPatterns {
		matches := regexp.MustCompile(pattern).FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				url := normalizeURL(match[len(match)-1], baseURL)
				if url != "" {
					urls = append(urls, url)
				}
			}
		}
	}
	return RemoveDuplicates(urls)
}

// normalizeURL 标准化URL
func normalizeURL(path, baseURL string) string {
	if strings.Contains(path, "http") {
		return ""
	}

	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "../", "/")
	if path == "" {
		return ""
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL + path
}
