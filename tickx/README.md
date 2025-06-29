### tickx

`tickx` 聚焦 **日期 / 时间** 相关操作，弥补 `time` 包的空缺。

---

## 安装
```bash
go get github.com/llyb120/yoya/tickx
```

---

## 能力一览
| 文件 | 代表函数 | 说明 |
| ---- | -------- | ---- |
| `date.go` | `Today` `Yesterday` `ParseLayout` 等 | 快速获取常用日期、解析格式化字符串 |

---

## 示例
```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/tickx"
)

func main() {
    fmt.Println("今天: ", tickx.Today())
    fmt.Println("昨天: ", tickx.Yesterday())

    t, _ := tickx.ParseLayout("2006-01-02 15:04:05", "2023-12-31 23:59:59")
    fmt.Println(tickx.FormatDateTime(t))
}
```

---

## 许可协议
本项目遵循 MIT License。 