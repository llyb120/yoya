package supx

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRecord(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  20,
	}
	record := NewRecord(person)
	record.Put("city", "New York")

	jsonBs, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonBs))

	tp := record.GetType()
	fmt.Println(tp, tp.String())

}

// Start Generation Here
func TestRecordUnmarshal(t *testing.T) {
	// 反序列化测试：测试主体数据和扩展字段解析
	jsonStr := `{"name":"Alice", "age":30, "city":"Los Angeles", "country":"USA"}`

	// 初始化一个空的 Record 实例
	rec := new(Record[Person])
	err := json.Unmarshal([]byte(jsonStr), rec)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 检查主体数据是否正确解析
	if rec.Data.Name != "Alice" || rec.Data.Age != 30 {
		t.Errorf("主体数据不匹配, 期望 {Name: \"Alice\", Age: 30}, 得到: %+v", rec.Data)
	}

	// 检查扩展字段 "city"
	if city, ok := rec.Ext["city"]; !ok {
		t.Error("扩展字段 'city' 未找到")
	} else if city != "Los Angeles" {
		t.Errorf("扩展字段 'city' 不匹配, 期望: Los Angeles, 得到: %v", city)
	}

	// 检查扩展字段 "country"
	if country, ok := rec.Ext["country"]; !ok {
		t.Error("扩展字段 'country' 未找到")
	} else if country != "USA" {
		t.Errorf("扩展字段 'country' 不匹配, 期望: USA, 得到: %v", country)
	}
}
