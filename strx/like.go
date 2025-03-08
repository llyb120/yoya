package strx

import (
	"regexp"
	"strings"
)

type LikeType int

const (
	Number LikeType = iota
)

func Like[T string | *string, K string | LikeType](str T, pattern K, extPatterns ...K) bool {
	var _str string
	switch s := any(str).(type) {
	case string:
		_str = strings.TrimSpace(s)
	case *string:
		if s == nil {
			return false
		}
		_str = strings.TrimSpace(*s)
		*s = _str
	}

	flag := like(_str, pattern)
	if flag {
		return true
	}
	for _, extPattern := range extPatterns {
		flag = like(_str, extPattern)
		if flag {
			return true
		}
	}
	return false
}

func like[T string | LikeType](str string, pattern T) bool {
	switch p := any(pattern).(type) {
	case string:
		modeMatch := strings.Contains(p, "*")
		if modeMatch {
			// 将通配符*替换为正则表达式.*
			regexPattern := strings.ReplaceAll(p, "*", ".*")
			// 使用正则表达式进行匹配
			return regexp.MustCompile(regexPattern).MatchString(str)
		}
		return strings.EqualFold(str, p)
	case LikeType:
		switch p {
		case Number:
			return isNumber(str)
		}
	}

	return false
}

func isNumber(str string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(str)
}
