package fingerScan

import (
	"github.com/kN6jq/fingerScan/internal/core"
	"github.com/kN6jq/fingerScan/internal/model"
)

// ScanConfig 扫描配置
type ScanConfig struct {
	ThreadCount int    // 扫描线程数
	OutputFile  string // 输出文件
	ProxyURL    string // 代理URL
}

// ScanResult 扫描结果
type ScanResult = model.ScanResult

// Scanner 指纹扫描器接口
type Scanner struct {
	scanner *core.Scanner
}

// NewScanner 创建新的扫描器实例
func NewScanner(urls []string, config ScanConfig) (*Scanner, error) {
	coreConfig := core.ScanConfig{
		ThreadCount: config.ThreadCount,
		OutputFile:  config.OutputFile,
		ProxyURL:    config.ProxyURL,
	}

	s, err := core.NewScanner(urls, coreConfig)
	if err != nil {
		return nil, err
	}

	return &Scanner{scanner: s}, nil
}

// Start 开始扫描
func (s *Scanner) Start() error {
	return s.scanner.Start()
}

// ScanSingleURL 扫描单个URL
func ScanSingleURL(url string, proxy string) (*ScanResult, error) {
	config := ScanConfig{
		ThreadCount: 1,
		ProxyURL:    proxy,
	}

	scanner, err := NewScanner([]string{url}, config)
	if err != nil {
		return nil, err
	}

	err = scanner.Start()
	if err != nil {
		return nil, err
	}

	// 返回结果
	return nil, nil // 这里需要修改返回实际的扫描结果
}

// LoadURLsFromFile 从文件加载URL列表
func LoadURLsFromFile(filename string) []string {
	return core.LoadURLsFromFile(filename)
}
