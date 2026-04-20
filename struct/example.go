package _struct

//my-converter/
//├── go.mod          // 依赖管理
//├── go.sum          // 依赖校验
//├── main.go         // 入口文件（通常很小）
//├── cmd/            // 存放不同的命令行子命令逻辑
//│   └── root.go
//├── internal/       // 内部业务逻辑，外部项目无法 import 这里面的包（Java 没有强制限制，Go 有）
//│   ├── parser/     // 解析器模块
//│   │   └── txt.go
//│   └── converter/  // 转换模块
//│       └── epub.go
//└── pkg/            // (可选) 愿意暴露给别人使用的工具代码
