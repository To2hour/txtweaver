package internal

// Importer 定义如何从文件读取并转为 Book 结构。
type Importer interface {
	Import(path string) (*Book, error)
}

// Exporter 定义如何将 Book 结构写入目标文件。
type Exporter interface {
	Export(b *Book, outputPath string) error
}
