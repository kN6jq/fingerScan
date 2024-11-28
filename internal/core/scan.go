// Package core 实现了指纹扫描的核心功能
package core

import (
	"github.com/kN6jq/fingerScan/internal/model"
	"github.com/kN6jq/fingerScan/internal/utils"
	"github.com/panjf2000/ants/v2"
	"strings"
	"sync"
)

// Scanner 指纹扫描器
type Scanner struct {
	urlQueue     *Queue
	httpClient   *HTTPClient
	fingerprints *model.FingerprintDB
	results      *ScanResults
	workerPool   *ants.Pool
	wg           sync.WaitGroup
	config       ScanConfig
}

// ScanConfig 扫描配置
type ScanConfig struct {
	ThreadCount int
	OutputFile  string
	ProxyURL    string
}

// ScanResults 扫描结果
type ScanResults struct {
	sync.Mutex
	All   []model.ScanResult
	Focus []model.ScanResult
}

// NewScanner 创建新的扫描器实例
func NewScanner(urls []string, config ScanConfig) (*Scanner, error) {
	fingerprints, err := LoadFingerprints()
	if err != nil {
		return nil, err
	}

	pool, err := ants.NewPool(config.ThreadCount)
	if err != nil {
		return nil, err
	}

	scanner := &Scanner{
		urlQueue:     NewQueue(),
		httpClient:   NewHTTPClient(config.ProxyURL),
		fingerprints: fingerprints,
		results:      &ScanResults{},
		workerPool:   pool,
		config:       config,
	}

	// 初始化URL队列
	for _, url := range urls {
		scanner.urlQueue.Push([]string{url, "0"})
	}

	return scanner, nil
}

// Start 开始扫描
func (s *Scanner) Start() error {
	defer s.workerPool.Release()

	// 提交扫描任务
	for s.urlQueue.Len() > 0 {
		s.wg.Add(1)
		err := s.workerPool.Submit(func() {
			defer s.wg.Done()
			s.scanWorker()
		})
		if err != nil {
			return err
		}
	}

	s.wg.Wait()

	// 输出结果
	s.outputResults()
	return nil
}

// scanWorker 扫描工作协程
func (s *Scanner) scanWorker() {
	for s.urlQueue.Len() > 0 {
		data := s.urlQueue.Pop()
		urls, ok := data.([]string)
		if !ok {
			continue
		}

		resp, err := s.httpClient.DoRequest(urls[0])
		if err != nil {
			continue
		}

		// 处理JS跳转
		if urls[1] == "0" {
			for _, jsURL := range resp.JSURLs {
				s.urlQueue.Push([]string{jsURL, "1"})
			}
		}

		// 识别CMS
		cms := s.identifyCMS(resp)
		result := model.ScanResult{
			URL:        resp.URL,
			CMS:        strings.Join(cms, ","),
			Server:     resp.Server,
			StatusCode: resp.StatusCode,
			Length:     resp.Length,
			Title:      resp.Title,
		}

		// 保存结果
		s.results.Lock()
		s.results.All = append(s.results.All, result)
		if len(cms) > 0 {
			s.results.Focus = append(s.results.Focus, result)
		}
		s.results.Unlock()

		// 输出扫描进度
		s.printProgress(result)
	}
}

// identifyCMS 识别CMS
func (s *Scanner) identifyCMS(resp *model.HTTPResponse) []string {
	var cms []string
	for _, fp := range s.fingerprints.Fingerprints {
		if s.matchFingerprint(fp, resp) {
			cms = append(cms, fp.CMS)
		}
	}
	return utils.RemoveDuplicates(cms)
}

// matchFingerprint 匹配指纹
func (s *Scanner) matchFingerprint(fp model.Fingerprint, resp *model.HTTPResponse) bool {
	var content string
	switch fp.Location {
	case "body":
		content = resp.Body
		if fp.Method == "faviconhash" && len(fp.Keywords) > 0 {
			return resp.FaviconHash == fp.Keywords[0]
		}
	case "header":
		content = utils.HeadersToString(resp.Headers)
	case "title":
		content = resp.Title
	default:
		return false
	}

	switch fp.Method {
	case "keyword":
		return utils.ContainsAllKeywords(content, fp.Keywords)
	case "regular":
		return utils.MatchesAllPatterns(content, fp.Keywords)
	}
	return false
}

// outputResults 输出扫描结果
func (s *Scanner) outputResults() {
	utils.PrintColoredResults(s.results.Focus)
	if s.config.OutputFile != "" {
		utils.SaveResults(s.config.OutputFile, s.results.All)
	}
}

// printProgress 打印扫描进度
func (s *Scanner) printProgress(result model.ScanResult) {
	if len(result.CMS) > 0 {
		utils.PrintColoredResult(result)
	} else {
		utils.PrintResult(result)
	}
}
