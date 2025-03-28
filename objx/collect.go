package objx

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goccy/go-reflect"
)

type selector struct {
	src string
	len int
	idx int
}

type selectorNode struct {
	key   string
	props []*selectorProp
}

var (
	opEqual = "=="
	opLike  = "*="
	opNot   = "!="
)

type selectorProp struct {
	key   string
	op    string
	value string
}

func (s *selectorNode) setProp(key string, value string) {
	value = strings.TrimPrefix(value, "=")
	if strings.HasSuffix(key, opLike) {
		s.props = append(s.props, &selectorProp{
			key:   key,
			op:    opLike,
			value: value,
		})
	} else if strings.HasSuffix(key, opNot) {
		s.props = append(s.props, &selectorProp{
			key:   key,
			op:    opNot,
			value: value,
		})
	} else {
		s.props = append(s.props, &selectorProp{
			key:   key,
			op:    opEqual,
			value: value,
		})
	}
}

func (s *selector) parse() []*selectorNode {
	s.src += " "
	s.len = len(s.src)
	s.idx = 0
	var buf bytes.Buffer
	var nodes []*selectorNode
	var current *selectorNode
	for s.idx < s.len {
		c := s.src[s.idx]
		if s.isWord(c) {
			if current == nil {
				current = &selectorNode{
					key: "",
				}
				nodes = append(nodes, current)
			}
			buf.WriteByte(c)
		} else if s.isSpace(c) {
			if current != nil {
				if buf.Len() > 0 {
					current.key = buf.String()
					buf.Reset()
				}
				current = nil
			}
		} else if c == '[' {
			s.idx++
			if current == nil {
				current = &selectorNode{
					key: "",
				}
				nodes = append(nodes, current)
			} else {
				if buf.Len() > 0 {
					current.key = buf.String()
					buf.Reset()
				}
			}
			// 读取到]
			for s.idx < s.len {
				if s.src[s.idx] != ']' {
					buf.WriteByte(s.src[s.idx])
					s.idx++
				} else {
					s.parseExpr(current, buf.String())
					buf.Reset()
					current = nil
					break
				}
			}
		}
		s.idx++
	}

	return nodes
}

func (s *selector) isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (s *selector) isWord(c byte) bool {
	return !s.isSpace(c) && c != '[' && c != ']'
}

func (s *selector) parseExpr(node *selectorNode, str string) {
	// return c >= '0' && c <= '9'
	// 解析形如 a=1,b=2,c="343214 12341" 的表达式
	var key, value bytes.Buffer
	var inQuote bool
	var quoteChar byte

	for i := 0; i < len(str); i++ {
		c := str[i]

		// 处理引号
		if c == '"' || c == '\'' {
			if !inQuote {
				inQuote = true
				quoteChar = c
				continue
			} else if quoteChar == c {
				inQuote = false
				node.setProp(key.String(), value.String())
				key.Reset()
				value.Reset()
				continue
			}
		}

		// 处理非引号状态
		if !inQuote {
			if c == '=' {
				// 切换到值的解析
				value.WriteByte(c)
				continue
			}
			if c == ',' {
				// 完成一个键值对
				if key.Len() > 0 && value.Len() > 0 {
					node.setProp(key.String(), value.String())
				}
				key.Reset()
				value.Reset()
				continue
			}

			// 根据是否在等号前后添加到key或value
			if value.Len() == 0 {
				key.WriteByte(c)
			} else {
				value.WriteByte(c)
			}
		} else {
			// 在引号内直接添加字符
			value.WriteByte(c)
		}
	}

	// 处理最后一个键值对
	if key.Len() > 0 && value.Len() > 0 {
		node.setProp(key.String(), value.String())
	}

	// var buf bytes.Buffer
	// for i := 0; i < len(str); i++ {
	// 	if s.isNumber(str[i]) {
	// 		buf.WriteByte(str[i])
	// 	}
	// }
	// return buf.String()
}

// 从任意对象中收集元素
func Pick(src any, rule string) []any {
	selector := &selector{
		src: rule,
	}
	nodes := selector.parse()

	var results []any
	var stack []any
	var stackMap = make(map[any]*int)
	var nodeMatches = make(map[any]map[int]bool)        // 使用s作为key记录节点匹配情况
	var nodePropMatches = make(map[any]map[string]bool) // 使用s作为key记录属性匹配情况

	// 添加一个辅助函数来检查是否构成完整路径
	checkFullMatch := func() bool {
		if len(stack) == 0 {
			return false
		}

		// 直接检查是否有完整的路径匹配
		// 从后向前搜索匹配
		curNode := len(nodes) - 1

		// 从后向前遍历堆栈
		for i := len(stack) - 1; i >= 0; i-- {
			s := stack[i]
			matches, exists := nodeMatches[s]

			// 如果当前节点匹配了当前选择器
			if exists && matches[curNode] {
				curNode--

				// 已经找到了所有选择器的匹配
				if curNode < 0 {
					return true
				}
			}
		}

		// 打印调试信息
		fmt.Printf("无法找到完整路径，当前匹配到第 %d 个选择器，总共 %d 个选择器\n",
			len(nodes)-curNode-1, len(nodes))

		return false
	}

	// 检查是否已经匹配了所有需要的属性
	checkAllPropsMatched := func(s any, node *selectorNode) bool {
		if len(node.props) == 0 {
			return true
		}

		propMap, exists := nodePropMatches[s]
		if !exists {
			return false
		}

		for _, prop := range node.props {
			if !propMap[prop.key] {
				return false
			}
		}
		return true
	}

	Walk(src, func(s, k, v any) any {
		pos := stackMap[s]
		if pos != nil {
			stack = stack[:*pos+1]
		}
		var newLevel bool
		if len(stack) > 0 {
			if stack[len(stack)-1] != s {
				newLevel = true
			}
		} else {
			newLevel = true
		}

		ptr := reflect.ValueOf(s).Elem().UnsafeAddr()
		fmt.Println(ptr)

		// itemPos := findStackPos(s)
		// if itemPos == -1 {
		// 	stack = stack[:itemPos]
		// }

		if newLevel {
			// 新层级，进行处理
			stack = append(stack, s)
			idx := len(stack) - 1
			stackMap[s] = &idx
			// 为当前节点创建匹配记录
			if nodeMatches[s] == nil {
				nodeMatches[s] = make(map[int]bool)
				nodePropMatches[s] = make(map[string]bool)
			}

			// defer func() {
			// 	stack = stack[:len(stack)-1]
			// 	// 不删除匹配记录，保留所有信息
			// }()
		}

		stack = append(stack, k)

		// 获取当前节点的键
		strKey, ok := k.(string)
		if !ok {
			return Unchanged
		}

		// 检查当前节点是否匹配任何选择器节点
		for i, node := range nodes {
			if strings.EqualFold(node.key, strKey) || node.key == "" {
				// 记录属性匹配情况
				matchProp(v, node, nodePropMatches[s])

				// 只有当所有属性都匹配时，才算这个节点匹配成功
				if checkAllPropsMatched(s, node) {
					nodeMatches[s][i] = true

					// 如果当前节点匹配了最后一个选择器并且形成了完整路径，收集结果
					if i == len(nodes)-1 && checkFullMatch() {
						results = append(results, v)
					}
				}
			}
		}

		return Unchanged
	})

	return results
}

// matchProp 检查属性是否匹配，并记录匹配情况
func matchProp(value any, node *selectorNode, matchInfo map[string]bool) {
	if len(node.props) == 0 {
		return
	}

	// 将值转换为字符串
	strValue := toString(value)

	// 检查每个属性
	for _, prop := range node.props {
		matched := false

		switch prop.op {
		case opEqual:
			matched = (strValue == prop.value)
		case opLike:
			matched = strings.Contains(strValue, prop.value)
		case opNot:
			matched = (strValue != prop.value)
		}

		if matched {
			matchInfo[prop.key] = true
		}
	}
}

// matchNodeProps 检查节点是否匹配选择器的所有属性
func matchNodeProps(value any, node *selectorNode) bool {
	if len(node.props) == 0 {
		return true
	}

	// 在Walk函数中，value是单个属性的值，而不是整个map或结构体
	// 所以我们需要直接检查这个值是否满足所有属性条件

	// 将值转换为字符串
	strValue := toString(value)

	// 检查是否满足所有属性条件
	allMatch := true
	for _, prop := range node.props {
		matched := false

		switch prop.op {
		case opEqual:
			matched = (strValue == prop.value)
		case opLike:
			matched = strings.Contains(strValue, prop.value)
		case opNot:
			matched = (strValue != prop.value)
		}

		if !matched {
			allMatch = false
			break
		}
	}

	return allMatch
}

// 将任意值转换为字符串
func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return ""
	}
}

// 从匹配节点中提取目标值
func extractValue(value any, target string) any {
	if target == "" {
		return value
	}

	m, ok := value.(map[string]any)
	if !ok {
		return nil
	}

	// 处理嵌套路径，如 "a.b.c"
	parts := strings.Split(target, ".")
	current := m

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}

		next, ok := current[part].(map[string]any)
		if !ok {
			return nil
		}
		current = next
	}

	return nil
}
