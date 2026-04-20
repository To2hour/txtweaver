package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadTxt 应该返回数据和错误，而不是在函数内部直接 panic
func readTxt(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		// Java 习惯抛异常，Go 习惯把错误返回给调用者处理
		return nil, err
	}
	return content, nil
}
func InitBook(path string) *Book {
	fileName := filepath.Base(path)
	split := strings.Split(fileName, ".")
	bookName := split[0]
	book := &Book{
		Title: bookName,
	}
	txt, err := readTxt(path)
	if err != nil {
		panic("初始化失败！")
	}
	parseContent(book, txt)
	fmt.Println("book.Title " + book.Title)
	fmt.Println("book.Author " + book.Author)
	for _, value := range book.Chapters {
		fmt.Println(value.Title)
		fmt.Println(value.Content)
	}
	return book
}
func bookToString(book *Book) {
	marshal, err := json.Marshal(book)
}
