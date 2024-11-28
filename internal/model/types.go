// Package model 定义了指纹扫描所需的数据结构
package model

// HTTPResponse 表示HTTP响应的结构体
type HTTPResponse struct {
	URL         string              // 请求URL
	Body        string              // 响应体
	Headers     map[string][]string // 响应头
	Server      string              // 服务器信息
	StatusCode  int                 // 状态码
	Length      int                 // 响应长度
	Title       string              // 网页标题
	JSURLs      []string            // JavaScript URL列表
	FaviconHash string              // favicon哈希值
}

// ScanResult 表示扫描结果的结构体
type ScanResult struct {
	URL        string `json:"url"`        // 目标URL
	CMS        string `json:"cms"`        // CMS类型
	Server     string `json:"server"`     // 服务器类型
	StatusCode int    `json:"statuscode"` // HTTP状态码
	Length     int    `json:"length"`     // 响应长度
	Title      string `json:"title"`      // 网页标题
}

// Fingerprint 表示CMS指纹特征
type Fingerprint struct {
	CMS      string   `json:"cms"`      // CMS名称
	Method   string   `json:"method"`   // 匹配方法
	Location string   `json:"location"` // 匹配位置
	Keywords []string `json:"keyword"`  // 关键字列表
}

// FingerprintDB 表示指纹数据库
type FingerprintDB struct {
	Fingerprints []Fingerprint `json:"fingerprint"` // 指纹列表
}
