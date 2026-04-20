package internal

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// parseContent 传入字节数组，然后初始化book
func parseContent(book *Book, data []byte) {
	// 1. 将 byte 转为 string 并按行读取
	// 为了节省内存，我们可以使用 bufio.Scanner
	reader := strings.NewReader(string(data))
	scanner := bufio.NewScanner(reader)

	// 2. 定义章节标题的正则（匹配：第xxx章、第xxx节、序章等）
	// 这个正则可以根据需求不断完善
	file, err := os.ReadFile("internal/regex.txt")
	if err != nil {
		panic("正则找不到")
	}
	regex := string(file)
	re := regexp.MustCompile(regex)

	var currentChapter *chapter

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 3. 匹配标题
		if re.MatchString(line) {
			// 如果匹配到了新章节，把之前的章节存入 Book，开始新章节
			currentChapter = &chapter{Title: line}
			book.Chapters = append(book.Chapters, currentChapter)
		} else {
			// 4. 如果没匹配到标题，说明是正文，追加到当前章节
			if currentChapter != nil {
				currentChapter.Content += line + "\n\n"
			} else {
				// 处理开头没有章节名的特殊情况
				currentChapter = &chapter{Title: "前言"}
				currentChapter.Content += line + "\n\n"
				book.Chapters = append(book.Chapters, currentChapter)
			}
		}
	}
}
