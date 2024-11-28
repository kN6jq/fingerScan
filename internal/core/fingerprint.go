package core

import (
	_ "embed"
	"encoding/json"
	"github.com/kN6jq/fingerScan/internal/model"
	"regexp"
	"strings"
)

//go:embed fingerprints/finger.json
var fingerprintData string

// LoadFingerprints 加载指纹库
func LoadFingerprints() (*model.FingerprintDB, error) {
	var db model.FingerprintDB
	if err := json.Unmarshal([]byte(fingerprintData), &db); err != nil {
		return nil, err
	}
	return &db, nil
}

// GetFingerprint 获取指定CMS的指纹
func GetFingerprint(db *model.FingerprintDB, cms string) []model.Fingerprint {
	var fingerprints []model.Fingerprint
	for _, fp := range db.Fingerprints {
		if fp.CMS == cms {
			fingerprints = append(fingerprints, fp)
		}
	}
	return fingerprints
}

// MatchFingerprint 匹配指纹是否符合目标
func MatchFingerprint(fp model.Fingerprint, target map[string]string) bool {
	content, ok := target[fp.Location]
	if !ok {
		return false
	}

	switch fp.Method {
	case "keyword":
		for _, keyword := range fp.Keywords {
			if !strings.Contains(content, keyword) {
				return false
			}
		}
		return true
	case "regular":
		for _, pattern := range fp.Keywords {
			matched, err := regexp.MatchString(pattern, content)
			if err != nil || !matched {
				return false
			}
		}
		return true
	}
	return false
}
