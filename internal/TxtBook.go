package internal

// 章节
type Chapter struct {
	Title   string
	Content string
}

// 书籍
type Book struct {
	Title    string
	Author   string
	Chapters []*Chapter
}

// Importer 定义如何从文件读取并转为 Book 结构
type Importer interface {
	Import(path string) (*Book, error)
}

// Exporter 定义如何将 Book 结构写入目标文件
type Exporter interface {
	Export(b *Book, outputPath string) error
}
