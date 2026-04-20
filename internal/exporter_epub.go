package internal

import (
	"fmt"
	"strings"

	"github.com/bmaupin/go-epub"
)

type epubExporter struct{}

func (ex *epubExporter) Export(b *Book, outputPath string) error {
	e := epub.NewEpub(b.Title)
	e.SetAuthor(b.Author)

	for _, ch := range b.Chapters {
		paras := strings.Split(ch.Content, "\n")
		htmlContent := "<h1>" + ch.Title + "</h1>"
		for _, p := range paras {
			p = strings.TrimSpace(p)
			if p != "" {
				htmlContent += "<p>" + p + "</p>"
			}
		}

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
