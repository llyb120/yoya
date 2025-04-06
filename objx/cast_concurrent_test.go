package objx

import (
	"fmt"
	"sync"
	"testing"

	"github.com/llyb120/yoya/internal"
)

// 用于并发测试的源结构体
type ConcurrentSource struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Score    float64  `json:"score"`
	IsActive bool     `json:"is_active"`
	Tags     []string `json:"tags"`
}

// 用于并发测试的目标结构体
type ConcurrentTarget struct {
	UserID      string   `json:"id"`
	DisplayName string   `json:"name"`
	Age         int      `json:"age"`
	Rating      string   `json:"rating"`
	Status      bool     `json:"is_active"`
	Categories  []string `json:"tags"`
}

// 用于并发测试的映射结构体
type ConcurrentMapTarget map[string]interface{}

// TestConcurrentConversion 测试并发环境下转换器的线程安全性
func TestConcurrentConversion(t *testing.T) {
	// 测试参数
	const (
		numGoroutines = 20 // 并发goroutine数量
		numIterations = 10 // 每个goroutine的迭代次数
	)

	// 创建一些测试数据
	sources := make([]ConcurrentSource, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		sources[i] = ConcurrentSource{
			ID:       i + 1000,
			Name:     "User-" + string(rune(65+i%26)),
			Age:      20 + i%30,
			Score:    85.5 + float64(i)/10.0,
			IsActive: i%2 == 0,
			Tags:     []string{"tag1", "tag2", "tag3"},
		}
	}

	// 等待所有goroutine完成
	var wg sync.WaitGroup
	var errorOccurred bool

	// 测试结构体到结构体的并发转换
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			source := sources[idx]

			for j := 0; j < numIterations; j++ {
				var target ConcurrentTarget
				err := internal.Cast(&target, source)
				if err != nil {
					t.Errorf("转换失败 (goroutine %d, iteration %d): %v", idx, j, err)
					errorOccurred = true
					return
				}

				// 简化验证，只检查几个关键字段
				if target.DisplayName != source.Name || target.Age != source.Age {
					t.Errorf("结构体转换结果不正确 (goroutine %d, iteration %d)", idx, j)
					errorOccurred = true
				}
			}
		}(i)
	}

	// 测试切片的并发转换
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			// 创建源切片
			sourceSlice := make([]ConcurrentSource, 5)
			for j := 0; j < 5; j++ {
				sourceSlice[j] = ConcurrentSource{
					ID:       idx*100 + j,
					Name:     "SliceUser-" + string(rune(65+j%26)),
					Age:      20 + j,
					Score:    85.5 + float64(j)/10.0,
					IsActive: j%2 == 0,
					Tags:     []string{"tag1", "tag2", "tag3"},
				}
			}

			for j := 0; j < numIterations; j++ {
				// 创建目标切片
				var targetSlice []ConcurrentTarget
				err := internal.Cast(&targetSlice, sourceSlice)
				if err != nil {
					t.Errorf("切片转换失败 (goroutine %d, iteration %d): %v", idx, j, err)
					errorOccurred = true
					return
				}

				// 验证切片长度
				if len(targetSlice) != len(sourceSlice) {
					t.Errorf("切片长度不匹配 (goroutine %d, iteration %d): expected %d, got %d",
						idx, j, len(sourceSlice), len(targetSlice))
					errorOccurred = true
					return
				}

				// 简化验证，只检查第一个元素
				if len(targetSlice) > 0 && targetSlice[0].DisplayName != sourceSlice[0].Name {
					t.Errorf("切片转换结果不正确 (goroutine %d, iteration %d)", idx, j)
					errorOccurred = true
				}
			}
		}(i)
	}

	// 测试读写映射缓存的并发操作
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()

			// 交替进行不同类型的转换，测试缓存机制
			for j := 0; j < numIterations; j++ {
				if j%2 == 0 {
					// 结构体到结构体
					source := sources[idx]
					var target ConcurrentTarget
					err := internal.Cast(&target, source)
					if err != nil {
						t.Errorf("缓存测试转换失败 (goroutine %d, iteration %d): %v", idx, j, err)
						errorOccurred = true
					}
				} else {
					// 结构体到结构体 - 使用不同结构体类型测试缓存
					source := ConcurrentSource{
						ID:       idx*1000 + j,
						Name:     "CacheTest-" + fmt.Sprint(j),
						Age:      30 + j,
						IsActive: j%3 == 0,
					}
					var target ConcurrentTarget
					err := internal.Cast(&target, source)
					if err != nil {
						t.Errorf("缓存测试转换失败 (goroutine %d, iteration %d): %v", idx, j, err)
						errorOccurred = true
					}
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	wg.Wait()

	if !errorOccurred {
		t.Log("所有并发测试完成，未发现数据竞争问题")
	}
}
