package internal

// 章节
type chapter struct {
	Title   string
	Content string
}

// 书籍
type Book struct {
	Title    string
	Author   string
	Chapters []*chapter
}
