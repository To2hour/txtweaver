package parser

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"txtweaver/internal"
	"txtweaver/internal/common"
)

type txtInput struct{}

//go:embed regex.txt
var regexPattern string

func (t txtInput) Import(path string) (book *internal.Book, err error) {
	defer func() {
		if err != nil {
			fmt.Printf("传入txt 失败，路径: %s, 原因: %v\n", path, err)
		} else {
			fmt.Println("txt转book成功: " + book.Title)
		}
	}()
	book = common.InitBook(path, book)
	content, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("读取失败: %w", err)
	}
	parseContent(book, content)
	defer content.Close()
	return book, nil
}

// parseContent 传入字节数组，然后初始化book
func parseContent(book *internal.Book, data io.Reader) {
	scanner := bufio.NewScanner(data)
	// 2. 定义章节标题的正则（匹配：第xxx章、第xxx节、序章等）
	re := regexp.MustCompile(regexPattern)

	var currentChapter *internal.Chapter

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 3. 匹配标题
		if re.MatchString(line) {
			// 如果匹配到了新章节，把之前的章节存入 Book，开始新章节
			currentChapter = &internal.Chapter{Title: line}
			book.Chapters = append(book.Chapters, currentChapter)
		} else {
			// 4. 如果没匹配到标题，说明是正文，追加到当前章节
			if currentChapter != nil {
				currentChapter.Content += line + "\n\n"
			} else {
				// 处理开头没有章节名的特殊情况
				currentChapter = &internal.Chapter{Title: "前言"}
				currentChapter.Content += line + "\n\n"
				book.Chapters = append(book.Chapters, currentChapter)
			}
		}
	}
}
