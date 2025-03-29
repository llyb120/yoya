package objx

import (
	"bytes"
	"strings"
)

type selector struct {
	src string
	len int
	idx int
}

type selectorNode struct {
	key   string
	props map[string]*selectorProp
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
		s.props[key] = &selectorProp{
			key:   key,
			op:    opLike,
			value: value,
		}
	} else if strings.HasSuffix(key, opNot) {
		s.props[key] = &selectorProp{
			key:   key,
			op:    opNot,
			value: value,
		}
	} else {
		s.props[key] = &selectorProp{
			key:   key,
			op:    opEqual,
			value: value,
		}
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
					key:   "",
					props: make(map[string]*selectorProp),
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
					key:   "",
					props: make(map[string]*selectorProp),
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
