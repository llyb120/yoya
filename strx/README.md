### strx

`strx` 提供了一些针对字符串的 **简洁扩展函数**，当前重点功能为类 SQL 的 `LIKE` 模糊匹配。

---

## 安装
```bash
go get github.com/llyb120/yoya/strx
```

---

## API
| 函数 | 说明 |
| ---- | ---- |
| `Like(str, pattern, extPatterns...)` | 判断 `str` 是否符合 `pattern` / 多 pattern，可使用 `*` 作为通配符 |
| `LikeType` | 预置匹配类型，目前只有 `strx.Number` —— 是否为纯数字 |

---

## 示例
```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/strx"
)

func main() {
    fmt.Println(strx.Like("config.yaml", "*.yaml")) // true
    fmt.Println(strx.Like("12345", strx.Number))    // true

    // 多模式匹配：满足其一即可
    ok := strx.Like("img.png", "*.jpg", "*.png")
    fmt.Println(ok) // true
}
```

---

## 许可协议
本项目遵循 MIT License。 