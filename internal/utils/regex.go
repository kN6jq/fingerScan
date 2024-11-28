// Package utils 提供正则和输出相关工具函数
package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ContainsAllKeywords 检查字符串是否包含所有关键字
func ContainsAllKeywords(content string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(content, keyword) {
			return false
		}
	}
	return true
}

// MatchesAllPatterns 检查字符串是否匹配所有正则表达式
func MatchesAllPatterns(content string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, content)
		if err != nil || !matched {
			return false
		}
	}
	return true
}

// ExtractFaviconPaths 提取favicon路径
func ExtractFaviconPaths(content string) []string {
	re := regexp.MustCompile(`href="(.*?favicon....)"`)
	matches := re.FindAllStringSubmatch(content, -1)

	var paths []string
	for _, match := range matches {
		if len(match) >= 2 {
			paths = append(paths, match[1])
		}
	}
	return paths
}

// HeadersToString 将HTTP头转换为字符串
func HeadersToString(headers map[string][]string) string {
	data, _ := json.Marshal(headers)
	return string(data)
}
