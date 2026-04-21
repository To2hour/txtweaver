package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"txtweaver/internal"

	"github.com/spf13/cobra"
)

var (
	inPath    string
	outPath   string
	outFormat string
	inFormat  string
	author    string
)

var rootCmd = &cobra.Command{
	Use:   "txtweaver",
	Short: "将文本等格式转换为电子书",
	Long: `txtweaver：命令行电子书转换工具。

在根命令上指定输入文件与输出格式即可转换；可选 -out 指定输出路径（默认与输入同目录、同主文件名）。
子命令 version 用于查看版本。`,
	RunE: runConvert,
}

func runConvert(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(inPath) == "" || strings.TrimSpace(outFormat) == "" {
		_ = cmd.Usage()
		return fmt.Errorf("必须提供 --in（或 -i）与 --out-format（或 -f）")
	}

	inFmt := strings.TrimSpace(inFormat)
	if inFmt == "" {
		inFmt = strings.TrimPrefix(strings.ToLower(filepath.Ext(inPath)), ".")
	}
	if inFmt == "" {
		inFmt = "txt"
	}

	outFmt := strings.TrimSpace(outFormat)

	dest := strings.TrimSpace(outPath)
	if dest == "" {
		ext := "." + strings.ToLower(strings.TrimPrefix(outFmt, "."))
		base := filepath.Base(inPath)
		if e := filepath.Ext(base); e != "" {
			base = strings.TrimSuffix(base, e)
		}
		dest = filepath.Join(filepath.Dir(inPath), base+ext)
	}

	imp, err := internal.GetImporter(inFmt)
	if err != nil {
		return err
	}
	exp, err := internal.GetExporter(outFmt)
	if err != nil {
		return err
	}

	book, err := imp.Import(inPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(author) != "" {
		book.Author = strings.TrimSpace(author)
	}

	if err := exp.Export(book, dest); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "导出成功: "+dest)
	return nil
}

func init() {
	f := rootCmd.Flags()
	f.StringVarP(&inPath, "in", "i", "", "输入文件路径")
	f.StringVarP(&outFormat, "out-format", "f", "", "输出格式（如 epub）")
	f.StringVarP(&outPath, "out", "o", "", "输出文件路径（可选，默认与输入同目录同名）")
	f.StringVar(&inFormat, "in-format", "", "输入格式（可选，默认从扩展名推断）")
	f.StringVar(&author, "author", "", "作者（可选）")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
