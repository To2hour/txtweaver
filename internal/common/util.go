package common

import (
	"path/filepath"
	"strings"
	"txtweaver/internal"
)

func InitBook(path string, book *internal.Book) *internal.Book {
	bookName := filepath.Base(path)
	if split := strings.Split(bookName, "."); len(split) > 0 {
		bookName = split[0]
	}
	book = &internal.Book{
		Title: bookName,
	}
	return book
}
