package internal

// Chapter 表示一本书的章节。
type Chapter struct {
	Title   string
	Content string
}

// Book 表示一本书的结构。
type Book struct {
	Title    string
	Author   string
	Chapters []*Chapter
}
