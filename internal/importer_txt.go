package internal

import (
	"bufio"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type txtImporter struct{}

// 章节标题正则（覆盖常见“第X章/节/卷”、序章/楔子/后记等）。
const txtChapterTitlePattern = `^(?:(?:第[\p{Han}0-9零一二三四五六七八九十百千万两]+[章节卷部篇回])|(?:序章|楔子|引子|前言|后记|尾声|番外)).*$`

func (t *txtImporter) Import(path string) (book *Book, err error) {
	defer func() {
		if err != nil {
			fmt.Printf("传入txt失败，路径: %s, 原因: %v\n", path, err)
		} else if book != nil {
			fmt.Println("txt转book成功: " + book.Title)
		}
	}()

	book = initBookFromPath(path)

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("读取失败: %w", err)
	}
	defer f.Close()

	parseTxtIntoBook(book, f)
	return book, nil
}
func parseTxtIntoBook(book *Book, f *os.File) {
	scanner := bufio.NewScanner(f)
	re := regexp.MustCompile(txtChapterTitlePattern)

	var currentChapter *Chapter

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if re.MatchString(line) {
			currentChapter = &Chapter{Title: line}
			book.Chapters = append(book.Chapters, currentChapter)
			continue
		}
		if currentChapter == nil {
			currentChapter = &Chapter{Title: "前言"}
			book.Chapters = append(book.Chapters, currentChapter)
		}
		// 统一存为 HTML 片段，避免导出时二次“猜格式”破坏渲染。
		currentChapter.Content += "<p>" + html.EscapeString(line) + "</p>\n"
	}
}
func initBookFromPath(path string) *Book {
	bookName := filepath.Base(path)
	if split := strings.Split(bookName, "."); len(split) > 0 {
		bookName = split[0]
	}
	return &Book{Title: bookName}
}
