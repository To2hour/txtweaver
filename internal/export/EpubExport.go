package export

import (
	"fmt"
	"strings"
	"txtweaver/internal"

	"github.com/bmaupin/go-epub"
)

type EpubExporter struct{}

func (ex *EpubExporter) Export(b *internal.Book, outputPath string) error {
	// 1. 创建一个新的 EPUB 实例，设置书名
	e := epub.NewEpub(b.Title)
	e.SetAuthor(b.Author)
	var _ internal.Exporter = (*EpubExporter)(nil)
	// 2. 遍历 Book 里的所有章节
	for _, ch := range b.Chapters {
		// EPUB 内部是 HTML 格式
		// 我们需要把 TXT 的换行符 \n 换成 HTML 的换行标签 <p>
		// 这样在阅读器里才会有分段效果
		paras := strings.Split(ch.Content, "\n")
		htmlContent := "<h1>" + ch.Title + "</h1>"
		for _, p := range paras {
			p = strings.TrimSpace(p)
			if p != "" {
				htmlContent += "<p>" + p + "</p>"
			}
		}

		// 3. 将章节添加到 EPUB 中
		_, err := e.AddSection(htmlContent, ch.Title, "", "")
		if err != nil {
			panic("export错误")
		}
	}

	// 4. 写出文件
	err := e.Write(outputPath)
	if err != nil {
		fmt.Println(err)
		panic("export错误")
	}
	return err
}
