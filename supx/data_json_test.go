package supx

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 测试用的结构体类型
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Age      int    `json:"age"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type Profile struct {
	Bio      string
	Location string `json:"location"`
	Website  string `json:"website,omitempty"`
}

type ComplexStruct struct {
	User    User              `json:"user"`
	Tags    []string          `json:"tags"`
	Meta    map[string]string `json:"meta"`
	Count   *int              `json:"count,omitempty"`
	Profile *Profile          `json:"profile,omitempty"`
}

func TestData_MarshalJSON_BasicStruct(t *testing.T) {
	data := NewData[User]()
	user := User{
		ID:    1,
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}
	active := true
	user.IsActive = &active

	data.Set(user)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 验证JSON包含预期的字段
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	if result["id"].(float64) != 1 {
		t.Errorf("期望id=1, 得到=%v", result["id"])
	}
	if result["name"].(string) != "张三" {
		t.Errorf("期望name=张三, 得到=%v", result["name"])
	}
	if result["email"].(string) != "zhangsan@example.com" {
		t.Errorf("期望email=zhangsan@example.com, 得到=%v", result["email"])
	}
	if result["is_active"].(bool) != true {
		t.Errorf("期望is_active=true, 得到=%v", result["is_active"])
	}
}

func TestData_MarshalJSON_WithExtraFields(t *testing.T) {
	data := NewData[User]()
	user := User{ID: 1, Name: "李四", Age: 30}
	data.Set(user)

	// 添加额外的map字段
	data["custom_field"] = "自定义值"
	data["number_field"] = 42
	data["bool_field"] = true

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 验证结构体字段
	if result["name"].(string) != "李四" {
		t.Errorf("期望name=李四, 得到=%v", result["name"])
	}

	// 验证额外字段
	if result["custom_field"].(string) != "自定义值" {
		t.Errorf("期望custom_field=自定义值, 得到=%v", result["custom_field"])
	}
	if result["number_field"].(float64) != 42 {
		t.Errorf("期望number_field=42, 得到=%v", result["number_field"])
	}
	if result["bool_field"].(bool) != true {
		t.Errorf("期望bool_field=true, 得到=%v", result["bool_field"])
	}
}

func TestData_MarshalJSON_NonStructType(t *testing.T) {
	data := NewData[string]()
	data.Set("测试字符串")
	data["extra"] = "额外字段"

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	if result[dataKey].(string) != "测试字符串" {
		t.Errorf("期望$data=测试字符串, 得到=%v", result[dataKey])
	}
	if result["extra"].(string) != "额外字段" {
		t.Errorf("期望extra=额外字段, 得到=%v", result["extra"])
	}
}

func TestData_MarshalJSON_EmptyData(t *testing.T) {
	data := NewData[User]()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("期望空对象, 得到=%v", result)
	}
}

func TestData_MarshalJSON_ComplexStruct(t *testing.T) {
	data := NewData[ComplexStruct]()
	count := 10
	complex := ComplexStruct{
		User:  User{ID: 1, Name: "王五", Age: 28},
		Tags:  []string{"开发者", "Go"},
		Meta:  map[string]string{"project": "yoya", "role": "admin"},
		Count: &count,
		Profile: &Profile{
			Bio:      "Go开发工程师",
			Location: "北京",
			Website:  "https://example.com",
		},
	}
	data.Set(complex)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	// 验证嵌套结构
	user := result["user"].(map[string]interface{})
	if user["name"].(string) != "王五" {
		t.Errorf("期望user.name=王五, 得到=%v", user["name"])
	}

	tags := result["tags"].([]interface{})
	if len(tags) != 2 || tags[0].(string) != "开发者" {
		t.Errorf("期望tags=[开发者, Go], 得到=%v", tags)
	}

	meta := result["meta"].(map[string]interface{})
	if meta["project"].(string) != "yoya" {
		t.Errorf("期望meta.project=yoya, 得到=%v", meta["project"])
	}
}

func TestData_UnmarshalJSON_BasicStruct(t *testing.T) {
	jsonStr := `{"id": 1, "name": "测试用户", "email": "test@example.com", "age": 25, "is_active": true}`

	var data Data[User]
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	user := data.Data()
	if user == nil {
		t.Fatal("期望非nil用户数据")
	}

	if user.ID != 1 {
		t.Errorf("期望ID=1, 得到=%d", user.ID)
	}
	if user.Name != "测试用户" {
		t.Errorf("期望Name=测试用户, 得到=%s", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Errorf("期望Email=test@example.com, 得到=%s", user.Email)
	}
	if user.IsActive == nil || !*user.IsActive {
		t.Errorf("期望IsActive=true, 得到=%v", user.IsActive)
	}
}

func TestData_UnmarshalJSON_WithExtraFields(t *testing.T) {
	jsonStr := `{"id": 2, "name": "用户2", "age": 30, "custom_field": "自定义", "number_field": 100}`

	var data Data[User]
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	user := data.Data()
	if user == nil {
		t.Fatal("期望非nil用户数据")
	}

	if user.Name != "用户2" {
		t.Errorf("期望Name=用户2, 得到=%s", user.Name)
	}

	// 验证额外字段
	if data["custom_field"].(string) != "自定义" {
		t.Errorf("期望custom_field=自定义, 得到=%v", data["custom_field"])
	}
	if data["number_field"].(float64) != 100 {
		t.Errorf("期望number_field=100, 得到=%v", data["number_field"])
	}
}

func TestData_UnmarshalJSON_NonStructType(t *testing.T) {
	jsonStr := `{"$data": "直接字符串", "extra": "额外数据"}`

	var data Data[string]
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	str := data.Data()
	if str == nil {
		t.Fatal("期望非nil字符串数据")
	}

	// 注意：由于实现中对非结构体类型是直接存储&dataValue，所以这里可能需要特殊处理
	if data["extra"].(string) != "额外数据" {
		t.Errorf("期望extra=额外数据, 得到=%v", data["extra"])
	}
}

func TestData_UnmarshalJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"id": 1, "name": "未闭合的JSON`

	var data Data[User]
	if err := json.Unmarshal([]byte(invalidJSON), &data); err == nil {
		t.Fatal("期望解析无效JSON时出错")
	}
}

func TestData_UnmarshalJSON_TypeMismatch(t *testing.T) {
	// 测试类型不匹配的情况
	jsonStr := `{"id": "应该是数字", "name": 123, "age": "应该是数字"}`

	var data Data[User]
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	user := data.Data()
	if user == nil {
		t.Fatal("期望非nil用户数据")
	}

	// 由于类型转换失败，字段应该保持零值
	if user.ID != 0 {
		t.Errorf("期望ID=0 (类型转换失败), 得到=%d", user.ID)
	}
}

func TestData_JSON_RoundTrip(t *testing.T) {
	// 测试完整的序列化->反序列化循环
	original := NewData[User]()
	user := User{
		ID:    99,
		Name:  "往返测试",
		Email: "roundtrip@example.com",
		Age:   35,
	}
	active := false
	user.IsActive = &active

	original.Set(user)
	original["custom"] = "自定义值"
	original["number"] = 42.5

	// 序列化
	jsonBytes, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 反序列化
	var restored Data[User]
	if err := json.Unmarshal(jsonBytes, &restored); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 验证数据完整性
	restoredUser := restored.Data()
	if restoredUser == nil {
		t.Fatal("期望非nil用户数据")
	}

	if !reflect.DeepEqual(*restoredUser, user) {
		t.Errorf("用户数据不匹配:\n期望: %+v\n得到: %+v", user, *restoredUser)
	}

	if restored["custom"].(string) != "自定义值" {
		t.Errorf("自定义字段不匹配: 期望=自定义值, 得到=%v", restored["custom"])
	}
	if restored["number"].(float64) != 42.5 {
		t.Errorf("数字字段不匹配: 期望=42.5, 得到=%v", restored["number"])
	}
}

func TestData_JSON_EmptyAndNil(t *testing.T) {
	// 测试空数据的序列化和反序列化
	empty := NewData[User]()

	jsonBytes, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("序列化空数据失败: %v", err)
	}

	var restored Data[User]
	if err := json.Unmarshal(jsonBytes, &restored); err != nil {
		t.Fatalf("反序列化空数据失败: %v", err)
	}

	if restored.Data() == nil {
		t.Errorf("期望非nil零值数据, 得到=nil")
	}
}

func TestData_JSON_PointerFields(t *testing.T) {
	// 测试包含指针字段的结构体
	data := NewData[ComplexStruct]()
	count := 100
	complex := ComplexStruct{
		User:  User{ID: 1, Name: "指针测试", Age: 25},
		Count: &count,
		Profile: &Profile{
			Bio:      "测试简介",
			Location: "上海",
		},
	}
	data.Set(complex)

	// 序列化
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 反序列化
	var restored Data[ComplexStruct]
	if err := json.Unmarshal(jsonBytes, &restored); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	restoredComplex := restored.Data()
	if restoredComplex == nil {
		t.Fatal("期望非nil复杂结构数据")
	}

	if restoredComplex.Count == nil || *restoredComplex.Count != 100 {
		t.Errorf("Count指针字段不匹配: 期望=100, 得到=%v", restoredComplex.Count)
	}

	if restoredComplex.Profile == nil {
		t.Fatal("期望非nil Profile")
	}
	if restoredComplex.Profile.Bio != "测试简介" {
		t.Errorf("Profile.Bio不匹配: 期望=测试简介, 得到=%s", restoredComplex.Profile.Bio)
	}
}

func BenchmarkData_MarshalJSON(b *testing.B) {
	data := NewData[User]()
	user := User{ID: 1, Name: "性能测试", Age: 30}
	data.Set(user)
	data["extra"] = "额外数据"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkData_UnmarshalJSON(b *testing.B) {
	jsonStr := `{"id": 1, "name": "性能测试", "age": 30, "extra": "额外数据"}`
	jsonBytes := []byte(jsonStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var data Data[User]
		if err := json.Unmarshal(jsonBytes, &data); err != nil {
			b.Fatal(err)
		}
	}
}
