package main_test

import (
	"testing"
	"txtweaver/internal/export"
	"txtweaver/internal/parser"
)

func TestName(t *testing.T) {
	book, _ := parser.InputAndParser("internal/test.epub")
	// 注意 E 后面的冒号
	outputPath := "E:\\study\\goLang\\code\\txtweaver\\internal\\test1.epub"
	export.ExportToEpub(book, outputPath)
}
