// Package main 提供程序入口
package main

import (
	"flag"
	"fmt"
	"github.com/kN6jq/fingerScan/internal/core"
	"github.com/kN6jq/fingerScan/internal/utils"
	"github.com/kN6jq/fingerScan/pkg/logger"
	"os"
	"time"
)

var (
	config = struct {
		file   string
		url    string
		output string
		thread int
		proxy  string
	}{}
)

func init() {
	flag.StringVar(&config.file, "f", "", "待识别的文件")
	flag.StringVar(&config.url, "u", "", "待识别的url")
	flag.StringVar(&config.output, "o", "", "保存的文件名(json或csv)")
	flag.IntVar(&config.thread, "t", 100, "扫描线程")
	flag.StringVar(&config.proxy, "p", "", "代理")
	flag.Parse()
}

func main() {
	startTime := time.Now()

	scanConfig := core.ScanConfig{
		ThreadCount: config.thread,
		OutputFile:  config.output,
		ProxyURL:    config.proxy,
	}

	var urls []string
	switch {
	case config.file != "":
		urls = utils.RemoveDuplicates(core.LoadURLsFromFile(config.file))
	case config.url != "":
		urls = []string{config.url}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	scanner, err := core.NewScanner(urls, scanConfig)
	if err != nil {
		logger.Error("初始化扫描器失败:", err)
		os.Exit(1)
	}

	if err := scanner.Start(); err != nil {
		logger.Error("扫描过程出错:", err)
		os.Exit(1)
	}

	fmt.Printf("扫描完成，耗时: %v\n", time.Since(startTime))
}
