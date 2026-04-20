package internal

import (
	"fmt"
	"strings"
)

// GetImporter 根据输入格式返回对应 Importer。
// 当前支持: "txt"
func GetImporter(format string) (Importer, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "txt":
		return &txtImporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported importer format: %s", format)
	}
}

// GetExporter 根据输出格式返回对应 Exporter。
// 当前支持: "epub"
func GetExporter(format string) (Exporter, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "epub":
		return &epubExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported exporter format: %s", format)
	}
}
