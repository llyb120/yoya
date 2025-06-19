package objx

import (
	"fmt"
	"testing"
	"time"

	"github.com/llyb120/yoya/syncx"
)

func TestWalk(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var item = map[string][]Person{
		"name": {{Name: "张三", Age: 30}, {Name: "李四", Age: 25}},
	}
	now := time.Now()
	Walk(item, func(s, k any, v any) any {
		if k == "Name" {
			var fn = syncx.Async_0_1(func() string {
				time.Sleep(2 * time.Second)
				return "ok"
			})
			res := fn()
			if err := syncx.Await(res); err != nil {
				return nil
			}
			return res
		}
		return nil
	}, Async, 1*Level)

	fmt.Println(item)
	fmt.Println(time.Since(now))
}

func TestBreakWalk(t *testing.T) {
	type Person struct {
		Name  string
		Age   int
		Child []Person
	}
	var item = map[string][]Person{
		"name": {{Name: "张三", Age: 30, Child: []Person{{Name: "张三儿子", Age: 30}, {Name: "张三女儿", Age: 25}}}, {Name: "李四", Age: 25}},
	}
	Walk(item, func(s, k any, v any) any {
		fmt.Println(k, v)
		return nil
	})

	println("---------------------------------------")

	Walk(item, func(s, k any, v any) any {
		if k == "Child" {
			return BreakWalkSelf
		}
		fmt.Println(k, v)
		return nil
	})
}

func TestEnsure(t *testing.T) {
	var a any
	var b any
	var c any
	a = 1
	b = 2
	c = 3

	var d int
	var e int
	var f int
	Ensure(&d, a, &e, b, &f, c)
	fmt.Println(d, e, f)
}

func TestCollect(t *testing.T) {
	src := "name[age=10,name*=张三] id [1,2,3] name[age=10,name='张三']"
	selector := &selector{
		src: src,
	}
	nodes := selector.parse()
	fmt.Printf("%+v\n", nodes)

	data := map[string]any{
		"name": []any{
			map[string]any{
				"age":  10,
				"name": "张三",
				"id":   1,
				"child": []map[string]any{
					{
						"child": 1,
					},
				},
			},
			map[string]any{
				"age":  10,
				"name": "张三",
				"id":   1,
			},
		},
	}
	results := Pick[string](data, "name [age=10,name='张三'] id", Distinct)
	fmt.Printf("%+v\n", results)

	reuslt2 := Pick[any](data, "child")
	fmt.Printf("%+v\n", reuslt2)
}

func TestComplexPick(t *testing.T) {
	// 创建一个复杂的嵌套数据结构
	complexData := map[string]any{
		"users": []map[string]any{
			{
				"id": 1,
				"profile": map[string]any{
					"name": "张三",
					"age":  28,
					"skills": []map[string]any{
						{"name": "编程", "level": 9},
						{"name": "设计", "level": 7},
					},
					"contact": map[string]any{
						"email": "zhangsan@example.com",
						"phone": "13800138000",
					},
				},
				"posts": []map[string]any{
					{
						"id":      101,
						"title":   "如何学习Go语言",
						"content": "Go语言是一门很棒的语言...",
						"tags":    []string{"Go", "编程", "学习"},
						"comments": []map[string]any{
							{"user": "李四", "content": "非常有用的文章"},
							{"user": "王五", "content": "谢谢分享"},
						},
					},
					{
						"id":      102,
						"title":   "数据结构基础",
						"content": "理解数据结构对编程很重要...",
						"tags":    []string{"数据结构", "编程", "基础"},
						"comments": []map[string]any{
							{"user": "赵六", "content": "讲解得很清楚"},
						},
					},
				},
			},
			{
				"id": 2,
				"profile": map[string]any{
					"name": "李四",
					"age":  32,
					"skills": []map[string]any{
						{"name": "管理", "level": 8},
						{"name": "编程", "level": 6},
					},
					"contact": map[string]any{
						"email": "lisi@example.com",
						"phone": "13900139000",
					},
				},
				"posts": []map[string]any{
					{
						"id":      201,
						"title":   "项目管理技巧",
						"content": "有效的项目管理需要...",
						"tags":    []string{"管理", "项目", "技巧"},
						"comments": []map[string]any{
							{"user": "张三", "content": "学到了很多"},
						},
					},
				},
			},
		},
		"categories": []map[string]any{
			{"id": 1, "name": "技术"},
			{"id": 2, "name": "管理"},
			{"id": 3, "name": "设计"},
		},
	}

	// 测试1：查找所有技能等级大于7的技能
	skills := Pick[map[string]any](complexData, "skills [level>7]")
	fmt.Println("高级技能:")
	for _, skill := range skills {
		fmt.Printf("  %s (等级: %v)\n", skill["name"], skill["level"])
	}

	// 测试2：查找所有张三的文章评论
	comments := Pick[map[string]any](complexData, "comments [user='张三']")
	fmt.Println("\n张三的评论:")
	for _, comment := range comments {
		fmt.Printf("  %s: %s\n", comment["user"], comment["content"])
	}

	// 测试3：查找所有带有"编程"标签的文章
	posts := Pick[map[string]any](complexData, "posts")
	fmt.Println("\n编程相关文章:")
	for _, post := range posts {
		tags, ok := post["tags"].([]string)
		if ok {
			for _, tag := range tags {
				if tag == "编程" {
					fmt.Printf("  %s\n", post["title"])
					break
				}
			}
		}
	}

	// 测试4：查找所有用户的联系方式
	contacts := Pick[map[string]any](complexData, "contact")
	fmt.Println("\n用户联系方式:")
	for i, contact := range contacts {
		fmt.Printf("  用户%d: 邮箱=%s, 电话=%s\n", i+1, contact["email"], contact["phone"])
	}

	// 测试5：使用多层嵌套选择器
	userPosts := Pick[map[string]any](complexData, "users profile[name='张三'] posts")
	fmt.Println("\n张三的所有文章:")
	for _, post := range userPosts {
		fmt.Printf("  %s\n", post["title"])
	}

	// 测试6：复杂条件组合
	result := Pick[map[string]any](complexData, "users [id=1] profile skills [level>5]")
	fmt.Println("\nID为1的用户的高级技能(等级>5):")
	for _, item := range result {
		fmt.Printf("  %s (等级: %v)\n", item["name"], item["level"])
	}

	// 测试7：使用Walk进行数据转换
	fmt.Println("\n将所有年龄增加1:")
	Walk(complexData, func(s any, k any, v any) any {
		if k == "age" {
			if age, ok := v.(int); ok {
				return age + 1
			}
		}
		return Unchanged
	})

	// 验证年龄是否已更新
	users := Pick[map[string]any](complexData, "profile")
	for _, user := range users {
		fmt.Printf("  %s: %d岁\n", user["name"], user["age"])
	}
}
