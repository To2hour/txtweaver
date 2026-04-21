# AGENTS.md

> 本文件面向 AI 编码助手（Cursor / Claude / Codex 等）。
> 项目体量还小，遵循以下约定可以让 Agent 快速、准确地协作。

## 项目简介

**txtweaver** 是一个用 Go 编写的命令行工具，用于把本地小说类 `txt` 文件转换为电子书格式（当前支持 `epub`，后续计划扩展 `mobi` / `pdf` 等）。

- 语言 / 版本：Go 1.26
- CLI 框架：[`spf13/cobra`](https://github.com/spf13/cobra)
- EPUB 生成：[`bmaupin/go-epub`](https://github.com/bmaupin/go-epub)
- 运行入口：`main.go` → `cmd.Execute()`

## 目录结构

```
txtweaver/
├── main.go                 # 极简入口，仅调用 cmd.Execute()
├── go.mod / go.sum         # 依赖管理
├── cmd/                    # Cobra 命令定义
│   ├── root.go             # 根命令
│   └── version.go          # version 子命令（示例 / 占位）
├── internal/               # 核心业务逻辑（对外不可 import）
│   ├── book.go             # 领域模型：Book / Chapter
│   ├── contract.go         # 接口：Importer / Exporter
│   ├── factory.go          # 根据格式字符串返回对应实现
│   ├── importer_txt.go     # txt → Book
│   ├── exporter_epub.go    # Book → epub
│   └── parser/
│       └── regex.txt       # 章节标题正则（通过 //go:embed 嵌入）
├── struct/
│   └── example.go          # 仅作结构说明的注释文件
└── AGENTS.md               # 本文件
```

## 架构约定

核心数据流非常直白：

```
txt file  ──▶  Importer  ──▶  *Book  ──▶  Exporter  ──▶  epub file
```

- `Book` / `Chapter` 是中间通用模型，新增任何输入/输出格式都应通过这个模型中转。
- 新增一种**输入格式**：
  1. 在 `internal/` 下新建 `importer_xxx.go`，实现 `Importer` 接口。
  2. 在 `internal/factory.go` 的 `GetImporter` switch 中注册。
- 新增一种**输出格式**：
  1. 在 `internal/` 下新建 `exporter_xxx.go`，实现 `Exporter` 接口。
  2. 在 `internal/factory.go` 的 `GetExporter` switch 中注册。
- 新增一个**CLI 子命令**：
  1. 在 `cmd/` 下新建 `xxx.go`，参考 `version.go` 的写法。
  2. 在该文件的 `init()` 里 `rootCmd.AddCommand(...)`。
- 章节识别规则统一写在 `internal/parser/regex.txt`，通过 `//go:embed` 注入，不要把规则硬编码到 Go 文件里。

## 编码规范

- 包名全部小写，目录名与包名保持一致（`struct/` 例外，其内部用了 `package _struct`，仅作示例）。
- 对外暴露的类型 / 函数用大驼峰；包内私有实现（如 `txtImporter`、`epubExporter`）保持小写开头。
- 错误处理：
  - 用 `fmt.Errorf("...: %w", err)` 做 wrap，保留上层可 `errors.Is/As`。
  - 面向用户的日志可用中文（项目现状如此），但 `error` 内容建议同时包含关键路径信息。
- 注释使用中文，符合现有风格；导出标识符保留 Go 官方风格的 `// TypeName ...` 开头。
- 不要引入额外的大型框架（ORM、Web 框架等）——本项目定位是单机 CLI。

## 常用命令

```bash
# 运行
go run .

# 构建（Windows 下生成 txtweaver.exe）
go build -o txtweaver.exe .

# 执行版本子命令
go run . version -T hello

# 格式化 & 静态检查
go fmt ./...
go vet ./...

# 跑测试（如有）
go test ./...
```

## Agent 工作守则

在对此仓库做改动时，请遵守：

1. **优先编辑已存在文件**，不要随意新建文件，尤其是 README / 文档类文件，除非用户明确要求。
2. **遵守 `internal/` 约定**：业务逻辑放 `internal/`，不要把解析 / 转换逻辑写进 `cmd/`。`cmd/` 只负责参数解析 + 调用 `internal`。
3. **新增格式时走工厂**：不要在 `cmd/` 里直接 `new` 出某个具体实现，要通过 `internal.GetImporter` / `GetExporter`。
4. **保持 `main.go` 极简**，仅做入口委托。
5. 修改章节识别相关内容时，改 `internal/parser/regex.txt`，而不是改 Go 代码里的字符串。
6. 构建产物（`*.exe`、`*.epub`、`dist/`、`output/` 等）已在 `.gitignore` 中，**不要提交**。
7. 修改代码后，若条件允许请跑一遍 `go build ./...` 与 `go vet ./...` 确认无编译错误。
8. 与用户沟通一律使用**简体中文**。

## 目前已知的 TODO / 粗糙点

- `cmd/root.go` 的 `Use` / `Short` / `Long` 还是模板文案（Hugo 字样），需要改成 txtweaver 自己的说明。
- `cmd/version.go` 只是示例，真实的 `convert` 命令（指定输入/输出路径与格式）尚未实现，后续应在这里扩展。
- `internal/factory.go` 中的 `Test` 全局变量只是调试残留，稳定前可以删掉。
- `struct/example.go` 仅为目录结构示意，不是真实代码，未来可移除或转为普通 markdown 文档。

---

如果后续项目结构发生较大变化（比如新增 `pkg/`、引入配置文件、接入网络下载等），请同步更新本文档。
