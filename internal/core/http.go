// Package core 实现了核心的HTTP请求和响应处理功能
package core

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/kN6jq/fingerScan/internal/model"
	"github.com/kN6jq/fingerScan/internal/utils"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

// HTTPClient 封装HTTP客户端功能
type HTTPClient struct {
	client *req.Client
	proxy  string
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(proxy string) *HTTPClient {
	client := req.C().
		EnableInsecureSkipVerify().
		SetUserAgent(utils.RandomUserAgent()).
		SetTLSFingerprintChrome().
		SetTimeout(defaultTimeout)

	if proxy != "" {
		client.SetProxyURL(proxy)
	}

	return &HTTPClient{
		client: client,
		proxy:  proxy,
	}
}

// DoRequest 执行HTTP请求
func (c *HTTPClient) DoRequest(urlStr string) (*model.HTTPResponse, error) {
	resp, err := c.client.R().Get(urlStr)
	if err != nil {
		// 尝试HTTP协议
		urlStr = strings.ReplaceAll(urlStr, "https://", "http://")
		resp, err = c.client.R().Get(urlStr)
		if err != nil {
			return nil, err
		}
	}

	body, err := resp.ToString()
	if err != nil {
		return nil, err
	}

	title := c.extractTitle(body)
	server := c.extractServer(resp.Header)
	faviconHash := c.getFaviconHash(body, urlStr)

	return &model.HTTPResponse{
		URL:         urlStr,
		Body:        body,
		Headers:     resp.Header,
		Server:      server,
		StatusCode:  resp.StatusCode,
		Length:      len(body),
		Title:       title,
		JSURLs:      utils.ExtractJSURLs(body, urlStr),
		FaviconHash: faviconHash,
	}, nil
}

// extractTitle 提取网页标题
func (c *HTTPClient) extractTitle(body string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	return strings.TrimSpace(strings.ReplaceAll(title, "\n", ""))
}

// extractServer 提取服务器信息
func (c *HTTPClient) extractServer(headers map[string][]string) string {
	if server := headers["Server"]; len(server) > 0 {
		return server[0]
	}
	if powered := headers["X-Powered-By"]; len(powered) > 0 {
		return powered[0]
	}
	return "None"
}

// getFaviconHash 获取favicon哈希值
func (c *HTTPClient) getFaviconHash(body, urlStr string) string {
	faviconURL := c.getFaviconURL(body, urlStr)
	return utils.CalculateFaviconHash(faviconURL)
}

// getFaviconURL 获取favicon URL
func (c *HTTPClient) getFaviconURL(body, urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	baseURL := u.Scheme + "://" + u.Host

	faviconPaths := utils.ExtractFaviconPaths(body)
	if len(faviconPaths) > 0 {
		fav := faviconPaths[0]
		switch {
		case strings.HasPrefix(fav, "//"):
			return "http:" + fav
		case strings.HasPrefix(fav, "http"):
			return fav
		default:
			return baseURL + "/" + strings.TrimPrefix(fav, "/")
		}
	}
	return baseURL + "/favicon.ico"
}
