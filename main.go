package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"txtweaver/internal"
)

func main() {
	inPath := flag.String("in", "", "输入文件路径（txt）")
	outPath := flag.String("out", "", "输出文件路径（epub）")
	author := flag.String("author", "", "作者（可选）")
	inFormat := flag.String("in-format", "", "输入格式（默认从扩展名推断，或 txt）")
	outFormat := flag.String("out-format", "", "输出格式（默认从扩展名推断，或 epub）")
	flag.Parse()

	if strings.TrimSpace(*inPath) == "" || strings.TrimSpace(*outPath) == "" {
		fmt.Fprintln(os.Stderr, "用法: txtweaver -in input.txt -out output.epub [-author name]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	inFmt := strings.TrimSpace(*inFormat)
	if inFmt == "" {
		inFmt = strings.TrimPrefix(strings.ToLower(filepath.Ext(*inPath)), ".")
	}
	if inFmt == "" {
		inFmt = "txt"
	}

	outFmt := strings.TrimSpace(*outFormat)
	if outFmt == "" {
		outFmt = strings.TrimPrefix(strings.ToLower(filepath.Ext(*outPath)), ".")
	}
	if outFmt == "" {
		outFmt = "epub"
	}

	imp, err := internal.GetImporter(inFmt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	exp, err := internal.GetExporter(outFmt)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	book, err := imp.Import(*inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if strings.TrimSpace(*author) != "" {
		book.Author = strings.TrimSpace(*author)
	}

	if err := exp.Export(book, *outPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("导出成功: " + *outPath)
}
