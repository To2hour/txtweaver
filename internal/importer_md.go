package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

type mdImporter struct{}

type heading struct {
	level    int
	lineIdx  int    // 0-based
	rawLine  string // 含 # 的原始行
	titleTxt string // 去掉 # 后的标题文本
}

func (m *mdImporter) Import(path string) (book *Book, err error) {
	defer func() {
		if err != nil {
			fmt.Printf("传入md失败，路径: %s, 原因: %v\n", path, err)
		} else if book != nil {
			fmt.Println("md转book成功: " + book.Title)
		}
	}()

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取失败: %w", err)
	}

	book = initBookFromPath(path)

	// 规则：
	// - 若第一个非空标题为 H1 且位于文件开头附近，则作为书名（Book.Title）
	// - 章节优先用 H2；若没有任何 H2，则用 H1（除去作为书名的那一个）
	// - 每章正文渲染为 HTML，尽可能保留 Markdown 展现效果
	lines := splitLines(string(src))
	lineStarts := computeLineStartOffsets(lines)

	heads := collectAtxHeadings(lines)
	if len(heads) == 0 {
		// 无标题：整篇作为一个章节
		bodyHTML := renderMarkdownToHTML(src)
		book.Chapters = []*Chapter{{Title: "正文", Content: bodyHTML}}
		return book, nil
	}

	// 尝试把文件顶部的首个 H1 当作书名（允许前面有空行）。
	titleHeadingLineLimit := 10
	titleHeadingIdx := -1
	for i, h := range heads {
		if h.level == 1 && h.lineIdx <= titleHeadingLineLimit {
			titleHeadingIdx = i
			if h.titleTxt != "" {
				book.Title = h.titleTxt
			}
			break
		}
	}

	hasH2 := false
	for _, h := range heads {
		if h.level == 2 {
			hasH2 = true
			break
		}
	}
	chapterLevel := 1
	if hasH2 {
		chapterLevel = 2
	}

	// 选出章节 headings（可能是 H2，或 H1）
	var chapterHeads []heading
	for i, h := range heads {
		if h.level != chapterLevel {
			continue
		}
		// 若用 H1 切章，并且某个 H1 被当作书名，则跳过它。
		if chapterLevel == 1 && i == titleHeadingIdx {
			continue
		}
		chapterHeads = append(chapterHeads, h)
	}

	// 如果没有章节标题（例如只有一个 H1 被当作书名），则整个正文为一个章节。
	if len(chapterHeads) == 0 {
		bodyStartLine := 0
		if titleHeadingIdx >= 0 {
			bodyStartLine = heads[titleHeadingIdx].lineIdx + 1
		}
		body := sliceByLineRange(src, lineStarts, bodyStartLine, len(lines))
		bodyHTML := renderMarkdownToHTML(body)
		book.Chapters = []*Chapter{{Title: "正文", Content: bodyHTML}}
		return book, nil
	}

	// 为每个章节取 markdown 片段：从该 heading 下一行到下一个同级 heading 行（不含）。
	for i, h := range chapterHeads {
		startLine := h.lineIdx + 1
		endLine := len(lines)
		if i+1 < len(chapterHeads) {
			endLine = chapterHeads[i+1].lineIdx
		}

		body := sliceByLineRange(src, lineStarts, startLine, endLine)
		bodyHTML := renderMarkdownToHTML(body)
		book.Chapters = append(book.Chapters, &Chapter{
			Title:   h.titleTxtOrFallback(),
			Content: bodyHTML,
		})
	}

	// 若没有书名，则沿用文件名（initBookFromPath 已做）。
	_ = filepath.Ext(path)
	return book, nil
}

func (h heading) titleTxtOrFallback() string {
	if strings.TrimSpace(h.titleTxt) != "" {
		return strings.TrimSpace(h.titleTxt)
	}
	return strings.TrimSpace(strings.TrimLeft(h.rawLine, "#"))
}

func renderMarkdownToHTML(src []byte) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			// 允许原始 HTML（对“尽可能保留效果”更友好；如果以后要安全输出可再收紧）
			gmhtml.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	_ = md.Convert(src, &buf)
	return buf.String()
}

func splitLines(s string) []string {
	// 统一为 \n 再 split，保留最后一行即使为空。
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

func computeLineStartOffsets(lines []string) []int {
	// lineStarts[i] 为第 i 行在原始文本中的 byte 起始偏移（以 \n 作为换行）。
	starts := make([]int, 0, len(lines)+1)
	offset := 0
	for _, line := range lines {
		starts = append(starts, offset)
		offset += len(line) + 1 // + '\n'
	}
	// 额外追加“末尾”，方便切片
	starts = append(starts, offset)
	return starts
}

func sliceByLineRange(src []byte, lineStarts []int, startLine, endLine int) []byte {
	if startLine < 0 {
		startLine = 0
	}
	if endLine < startLine {
		endLine = startLine
	}
	if startLine >= len(lineStarts) {
		return nil
	}
	if endLine >= len(lineStarts) {
		endLine = len(lineStarts) - 1
	}

	start := lineStarts[startLine]
	end := lineStarts[endLine]
	if start < 0 {
		start = 0
	}
	if end > len(src) {
		end = len(src)
	}
	if start > end {
		start = end
	}
	return src[start:end]
}

func collectAtxHeadings(lines []string) []heading {
	// 简化实现：只识别 ATX headings（#...######）。
	// 依赖 goldmark 渲染保留效果；切章本身先满足主流 md 写法。
	var hs []heading

	inFence := false
	for i, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if !strings.HasPrefix(line, "#") {
			continue
		}

		level := 0
		for level < len(line) && level < 6 && line[level] == '#' {
			level++
		}
		if level == 0 || level > 6 {
			continue
		}
		// 需要至少一个空格才当作 heading（避免 ######foo 误判）
		if len(line) > level && line[level] != ' ' && line[level] != '\t' {
			continue
		}
		title := strings.TrimSpace(line[level:])
		hs = append(hs, heading{
			level:    level,
			lineIdx:  i,
			rawLine:  raw,
			titleTxt: title,
		})
	}
	return hs
}
