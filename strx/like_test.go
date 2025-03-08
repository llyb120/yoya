package strx

import (
	"fmt"
	"testing"
)

func TestLike(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		pattern  any
		expected bool
	}{
		{"完全匹配", "hello ", "hello", true},
		{"大小写不敏感", "HELLO ", "hello", true},
		{"通配符匹配", "hello world", "hello*", true},
		{"通配符不匹配", "hi world", "hello*", false},
		{"数字匹配", "123", Number, true},
		{"非数字匹配", "abc", Number, false},
		{"空字符串", "", "", true},
		{"nil指针", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pattern any
			if tt.name == "nil指针" {
				pattern = nil
			} else if tt.name == "数字匹配" || tt.name == "非数字匹配" {
				pattern = Number
			} else {
				pattern = tt.pattern
			}

			var str any = tt.str
			if tt.name == "nil指针" {
				var s *string
				str = s
			}

			var result bool
			switch s := any(pattern).(type) {
			case string:
				result = Like(str.(string), s)
			case LikeType:
				result = Like(str.(string), s)
			}
			if result != tt.expected {
				t.Errorf("Like(%v, %v) = %v, want %v", tt.str, pattern, result, tt.expected)
			}
		})
	}

	var v = " foo "
	if Like(&v, "a", "foo") {
		fmt.Println("**", v, "**")
	}
	var v2 = "  "
	if Like(v2, "") {
		fmt.Println("**", v2, "**")
	}
}
