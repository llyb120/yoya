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
	opGt    = ">"
	opGe    = ">="
	opLt    = "<"
	opLe    = "<="
	opErr   = "err"
)

type selectorProp struct {
	key   string
	op    string
	value any
}

func (s *selectorNode) setProp(key string, value string) {
	if strings.HasPrefix(key, opLike) {
		s.props[key] = &selectorProp{
			key:   key,
			op:    opLike,
			value: value[2:],
		}
	} else if strings.HasPrefix(key, opNot) {
		s.props[key] = &selectorProp{
			key:   key,
			op:    opNot,
			value: value[2:],
		}
	} else if strings.HasPrefix(value, opGe) {
		val, ok := toFloat64(value[2:])
		s.props[key] = &selectorProp{
			key:   key,
			op:    opGe,
			value: val,
		}
		if !ok {
			s.props[key].op = opErr
		}
	} else if strings.HasPrefix(value, opGt) {
		val, ok := toFloat64(value[1:])
		s.props[key] = &selectorProp{
			key:   key,
			op:    opGt,
			value: val,
		}
		if !ok {
			s.props[key].op = opErr
		}
	} else if strings.HasPrefix(value, opLe) {
		val, ok := toFloat64(value[2:])
		s.props[key] = &selectorProp{
			key:   key,
			op:    opLe,
			value: val,
		}
		if !ok {
			s.props[key].op = opErr
		}
	} else if strings.HasPrefix(value, opLt) {
		val, ok := toFloat64(value[1:])
		s.props[key] = &selectorProp{
			key:   key,
			op:    opLt,
			value: val,
		}
		if !ok {
			s.props[key].op = opErr
		}
	} else {
		s.props[key] = &selectorProp{
			key:   key,
			op:    opEqual,
			value: value[1:],
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

func (s *selector) next(str string, pos int) byte {
	if pos >= len(str) {
		return 0
	}
	return str[pos+1]
}

func (s *selector) parseExpr(node *selectorNode, str string) {
	// return c >= '0' && c <= '9'
	// 解析形如 a=1,b=2,c="343214 12341" 的表达式
	var key, value bytes.Buffer
	var inQuote bool
	var quoteChar byte

	var i int
	for ; i < len(str); i++ {
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
			if c == '>' || c == '<' {
				value.WriteByte(c)
				if s.next(str, i) == '=' {
					value.WriteByte('=')
				}
				continue
			}
			if c == '!' || c == '*' {
				value.WriteByte(c)
				if s.next(str, i) == '=' {
					value.WriteByte('=')
				}
				continue
			}
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
