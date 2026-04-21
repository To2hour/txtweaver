package internal

import (
	"fmt"

	"github.com/bmaupin/go-epub"
)

type epubExporter struct{}

func (ex *epubExporter) Export(b *Book, outputPath string) error {
	e := epub.NewEpub(b.Title)
	e.SetAuthor(b.Author)

	for _, ch := range b.Chapters {
		// Content 约定为“章节正文的 HTML 片段”（由 Importer 负责产出），Exporter 只负责包装章节标题。
		htmlContent := "<h1>" + ch.Title + "</h1>\n" + ch.Content

		_, err := e.AddSection(htmlContent, ch.Title, "", "")
		if err != nil {
			return fmt.Errorf("添加章节失败(%s): %w", ch.Title, err)
		}
	}

	if err := e.Write(outputPath); err != nil {
		return fmt.Errorf("写出epub失败: %w", err)
	}
	return nil
}
