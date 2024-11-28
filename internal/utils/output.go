package utils

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gookit/color"
	"github.com/kN6jq/fingerScan/internal/model"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// PrintResult 打印普通扫描结果
func PrintResult(result model.ScanResult) {
	format := "[ %s | %s | %s | %d | %d | %s ]\n"
	fmt.Printf(format,
		result.URL,
		result.CMS,
		result.Server,
		result.StatusCode,
		result.Length,
		result.Title,
	)
}

// PrintColoredResult 打印带颜色的扫描结果
func PrintColoredResult(result model.ScanResult) {
	format := "[ %s | %s | %s | %d | %d | %s ]\n"
	color.RGBStyleFromString("237,64,35").Printf(format,
		result.URL,
		result.CMS,
		result.Server,
		result.StatusCode,
		result.Length,
		result.Title,
	)
}

// PrintColoredResults 打印重点资产结果
func PrintColoredResults(results []model.ScanResult) {
	color.RGBStyleFromString("244,211,49").Println("\n重点资产：")
	for _, result := range results {
		PrintColoredResult(result)
	}
}

// SaveResults 保存扫描结果
func SaveResults(filename string, results []model.ScanResult) error {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return SaveJSON(filename, results)
	case ".xlsx":
		return SaveXLSX(filename, results)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

// SaveJSON 保存JSON格式结果
func SaveJSON(filename string, results []model.ScanResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// SaveXLSX 保存XLSX格式结果
func SaveXLSX(filename string, results []model.ScanResult) error {
	xlsx := excelize.NewFile()
	headers := []string{"url", "cms", "server", "statuscode", "length", "title"}

	for i, header := range headers {
		col := string(rune('A' + i))
		xlsx.SetCellValue("Sheet1", col+"1", header)
	}

	for i, result := range results {
		row := strconv.Itoa(i + 2)
		xlsx.SetCellValue("Sheet1", "A"+row, result.URL)
		xlsx.SetCellValue("Sheet1", "B"+row, result.CMS)
		xlsx.SetCellValue("Sheet1", "C"+row, result.Server)
		xlsx.SetCellValue("Sheet1", "D"+row, result.StatusCode)
		xlsx.SetCellValue("Sheet1", "E"+row, result.Length)
		xlsx.SetCellValue("Sheet1", "F"+row, result.Title)
	}

	return xlsx.SaveAs(filename)
}
