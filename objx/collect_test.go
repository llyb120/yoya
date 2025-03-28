package objx

import (
	"fmt"
	"testing"
)

func TestPathMatch(t *testing.T) {
	// 测试用例1：简单路径
	t.Run("简单路径", func(t *testing.T) {
		data := map[string]any{
			"user": map[string]any{
				"name": "张三",
				"age":  30,
				"info": map[string]any{
					"address": "北京",
					"phone":   "123456",
				},
			},
		}

		// 测试: user info address
		results := Pick(data, "user info address")
		if len(results) != 1 {
			t.Errorf("期望匹配1个结果，实际匹配%d个", len(results))
		} else if results[0] != "北京" {
			t.Errorf("匹配结果错误，期望'北京'，实际是'%v'", results[0])
		}
	})

	// 测试用例2：带属性过滤的路径
	t.Run("带属性过滤的路径", func(t *testing.T) {
		data := map[string]any{
			"users": []any{
				map[string]any{
					"name": "张三",
					"age":  30,
					"id":   1,
				},
				map[string]any{
					"name": "李四",
					"age":  25,
					"id":   2,
				},
			},
		}

		// 测试: users[name="张三"] id
		results := Pick(data, `users[name="张三"] id`)
		fmt.Printf("测试结果: %v\n", results)
		if len(results) != 1 {
			t.Errorf("期望匹配1个结果，实际匹配%d个", len(results))
		} else if results[0] != 1 {
			t.Errorf("匹配结果错误，期望1，实际是'%v'", results[0])
		}
	})

	// 测试用例3：复杂嵌套结构
	t.Run("复杂嵌套结构", func(t *testing.T) {
		data := map[string]any{
			"departments": []any{
				map[string]any{
					"name": "技术部",
					"employees": []any{
						map[string]any{
							"name": "张三",
							"skills": []any{
								"Go", "Python", "JavaScript",
							},
						},
						map[string]any{
							"name": "李四",
							"skills": []any{
								"Java", "C++",
							},
						},
					},
				},
				map[string]any{
					"name": "市场部",
					"employees": []any{
						map[string]any{
							"name": "王五",
							"skills": []any{
								"Marketing", "Sales",
							},
						},
					},
				},
			},
		}

		// 测试: departments[name="技术部"] employees[name="张三"] skills
		results := Pick(data, `departments[name="技术部"] employees[name="张三"] skills`)
		fmt.Printf("测试结果: %v\n", results)
		if len(results) != 3 {
			t.Errorf("期望匹配3个结果，实际匹配%d个", len(results))
		}
	})
}
