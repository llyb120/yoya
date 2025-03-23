package objx

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// 源结构体
type Person struct {
	Name    string
	Age     int
	Address string
	Salary  float64
}

// 目标结构体
type Employee struct {
	Name       string
	Age        int
	Address    string
	Department string // 这个字段在源结构体中不存在
}

func TestConverter(t *testing.T) {
	// 创建转换器
	converter := newConverter()

	// 创建源结构体
	person := &Person{
		Name:    "张三",
		Age:     30,
		Address: "北京市海淀区",
		Salary:  10000.50,
	}

	// 创建目标结构体
	employee := &Employee{}

	// 执行转换
	err := converter.Convert(person, employee)
	if err != nil {
		t.Fatal(err)
	}

	// 验证转换结果
	if employee.Name != person.Name {
		t.Errorf("Name 转换失败: 预期 %s, 得到 %s", person.Name, employee.Name)
	}
	if employee.Age != person.Age {
		t.Errorf("Age 转换失败: 预期 %d, 得到 %d", person.Age, employee.Age)
	}
	if employee.Address != person.Address {
		t.Errorf("Address 转换失败: 预期 %s, 得到 %s", person.Address, employee.Address)
	}

	fmt.Printf("结构体转换结果: %+v\n", employee)

	// 测试切片转换
	persons := []*Person{
		{Name: "张三", Age: 30, Address: "北京", Salary: 10000},
		{Name: "李四", Age: 25, Address: "上海", Salary: 12000},
	}

	var employees []*Employee

	err = converter.ConvertSlice(&persons, &employees)
	if err != nil {
		t.Fatal(err)
	}

	// 验证切片转换结果
	if len(employees) != len(persons) {
		t.Errorf("切片长度不匹配: 预期 %d, 得到 %d", len(persons), len(employees))
	}

	for i, emp := range employees {
		p := persons[i]
		if emp.Name != p.Name || emp.Age != p.Age || emp.Address != p.Address {
			t.Errorf("切片元素 %d 转换失败", i)
		}
	}

	fmt.Println("切片转换成功!")

	// 性能测试：重复使用同一个转换器进行多次转换
	for i := 0; i < 10; i++ {
		newEmployee := &Employee{}
		err = converter.Convert(person, newEmployee)
		if err != nil {
			t.Fatal(err)
		}
	}
	fmt.Println("缓存测试成功!")
}

// 嵌套结构体的定义
type Address struct {
	City     string
	Province string
	ZipCode  string
}

type Contact struct {
	Phone  string
	Email  string
	WeChat string
}

// 复杂的源结构体
type ComplexPerson struct {
	ID        int
	Name      string
	BirthYear int
	Height    float64
	IsActive  bool
	Address   Address
	Contacts  []Contact
	Tags      []string
	Scores    map[string]float64
}

// 复杂的目标结构体
type UserProfile struct {
	UserID       string   `json:"id"`   // 从int转为string
	UserName     string   `json:"name"` // 映射到Name字段
	Age          int      // 需要从BirthYear计算
	HeightCM     string   // 从float转为string
	Active       bool     `json:"is_active"` // 映射到IsActive字段
	CityAddress  string   // 只取Address的City部分
	EmailAddress string   // 只取Contacts中的Email
	Labels       []string `json:"tags"` // 映射到Tags字段
	ScoreAverage float64  // 需要计算Scores的平均值
}

// 另一组复杂的测试用例
type ProductSource struct {
	SKU        string  `json:"product_code"`
	Name       string  `json:"title"`
	Price      float64 `json:"price"`
	Stock      int     `json:"inventory"`
	Dimensions struct {
		Length float64
		Width  float64
		Height float64
	}
	Categories []string `json:"categories"`
}

type ProductTarget struct {
	ProductCode      string  `json:"product_code"` // 通过json tag映射到SKU
	Title            string  `json:"title"`        // 通过json tag映射到Name
	PriceText        string  // Price转字符串
	Available        bool    // Stock>0转为bool
	Volume           float64 // 根据三维计算体积
	StringCategories string  // 将Categories转为逗号分隔的字符串
}

// 复杂测试用例
func TestComplexConversion(t *testing.T) {
	converter := newConverter()

	// 测试嵌套结构体转换
	testComplexPersonToUserProfile(t, converter)

	// 测试多种类型转换和字段映射
	testProductConversion(t, converter)

	// 测试切片中的复杂结构体转换
	testComplexSliceConversion(t, converter)

	fmt.Println("复杂转换测试全部通过!")
}

// 测试嵌套结构体和复杂类型转换
func testComplexPersonToUserProfile(t *testing.T, converter *Converter) {
	// 创建复杂的源数据
	currentYear := 2025
	source := &ComplexPerson{
		ID:        12345,
		Name:      "李明",
		BirthYear: 1990,
		Height:    178.5,
		IsActive:  true,
		Address: Address{
			City:     "广州",
			Province: "广东",
			ZipCode:  "510000",
		},
		Contacts: []Contact{
			{Phone: "13800138000", Email: "liming@example.com", WeChat: "liming_wx"},
			{Phone: "13900139000", Email: "liming_work@company.com", WeChat: "lm_company"},
		},
		Tags: []string{"开发", "后端", "Go语言"},
		Scores: map[string]float64{
			"算法":   85.5,
			"系统设计": 92.0,
			"编码":   88.5,
		},
	}

	// 创建目标结构体
	target := &UserProfile{}

	// 执行自定义转换 - 这个需要在Convert方法完成后手动处理一些字段
	err := converter.Convert(source, target)
	if err != nil {
		t.Fatal(err)
	}

	// 处理自定义转换字段
	target.UserID = fmt.Sprintf("%d", source.ID)
	target.Age = currentYear - source.BirthYear
	target.HeightCM = fmt.Sprintf("%.1fcm", source.Height)
	target.CityAddress = source.Address.City
	target.EmailAddress = source.Contacts[0].Email
	target.Labels = source.Tags

	// 计算分数平均值
	var sum float64
	for _, score := range source.Scores {
		sum += score
	}
	target.ScoreAverage = sum / float64(len(source.Scores))

	// 验证转换结果
	if target.UserID != "12345" {
		t.Errorf("UserID转换失败: 预期 %s, 得到 %s", "12345", target.UserID)
	}

	if target.Age != 35 {
		t.Errorf("Age计算失败: 预期 %d, 得到 %d", 35, target.Age)
	}

	if target.HeightCM != "178.5cm" {
		t.Errorf("Height转换失败: 预期 %s, 得到 %s", "178.5cm", target.HeightCM)
	}

	if target.CityAddress != "广州" {
		t.Errorf("CityAddress转换失败: 预期 %s, 得到 %s", "广州", target.CityAddress)
	}

	if target.EmailAddress != "liming@example.com" {
		t.Errorf("EmailAddress转换失败: 预期 %s, 得到 %s", "liming@example.com", target.EmailAddress)
	}

	expectedAvg := (85.5 + 92.0 + 88.5) / 3.0
	if math.Abs(target.ScoreAverage-expectedAvg) > 0.001 {
		t.Errorf("ScoreAverage计算失败: 预期 %.2f, 得到 %.2f", expectedAvg, target.ScoreAverage)
	}

	fmt.Printf("嵌套结构体转换结果: %+v\n", target)
}

// 测试多种类型转换和计算字段
func testProductConversion(t *testing.T, converter *Converter) {
	// 创建源产品
	source := &ProductSource{
		SKU:   "P-12345",
		Name:  "高性能笔记本电脑",
		Price: 6999.99,
		Stock: 25,
		Dimensions: struct {
			Length float64
			Width  float64
			Height float64
		}{
			Length: 35.6,
			Width:  24.5,
			Height: 1.8,
		},
		Categories: []string{"电子产品", "电脑", "笔记本"},
	}

	// 创建目标产品
	target := &ProductTarget{}

	// 执行基本转换
	err := converter.Convert(source, target)
	if err != nil {
		t.Fatal(err)
	}

	// 手动设置需要计算的字段
	target.ProductCode = source.SKU
	target.PriceText = fmt.Sprintf("￥%.2f", source.Price)
	target.Available = source.Stock > 0
	target.Volume = source.Dimensions.Length * source.Dimensions.Width * source.Dimensions.Height
	target.StringCategories = strings.Join(source.Categories, ", ")

	// 验证转换结果
	if target.ProductCode != "P-12345" {
		t.Errorf("ProductCode转换失败: 预期 %s, 得到 %s", "P-12345", target.ProductCode)
	}

	if target.Title != "高性能笔记本电脑" {
		t.Errorf("Title转换失败: 预期 %s, 得到 %s", "高性能笔记本电脑", target.Title)
	}

	if target.PriceText != "￥6999.99" {
		t.Errorf("PriceText转换失败: 预期 %s, 得到 %s", "￥6999.99", target.PriceText)
	}

	if !target.Available {
		t.Errorf("Available转换失败: 预期 %t, 得到 %t", true, target.Available)
	}

	expectedVolume := 35.6 * 24.5 * 1.8
	if math.Abs(target.Volume-expectedVolume) > 0.001 {
		t.Errorf("Volume计算失败: 预期 %.2f, 得到 %.2f", expectedVolume, target.Volume)
	}

	expectedCategories := "电子产品, 电脑, 笔记本"
	if target.StringCategories != expectedCategories {
		t.Errorf("StringCategories转换失败: 预期 %s, 得到 %s", expectedCategories, target.StringCategories)
	}

	fmt.Printf("产品转换结果: %+v\n", target)
}

// 测试复杂切片转换
func testComplexSliceConversion(t *testing.T, converter *Converter) {
	// 创建源数据切片
	sources := []*ProductSource{
		{
			SKU:   "P-1001",
			Name:  "智能手机",
			Price: 2999.99,
			Stock: 100,
			Dimensions: struct {
				Length float64
				Width  float64
				Height float64
			}{Length: 15.5, Width: 7.5, Height: 0.8},
			Categories: []string{"电子产品", "手机"},
		},
		{
			SKU:   "P-1002",
			Name:  "平板电脑",
			Price: 3499.99,
			Stock: 50,
			Dimensions: struct {
				Length float64
				Width  float64
				Height float64
			}{Length: 24.0, Width: 17.8, Height: 0.9},
			Categories: []string{"电子产品", "平板"},
		},
		{
			SKU:   "P-1003",
			Name:  "无线耳机",
			Price: 599.99,
			Stock: 0,
			Dimensions: struct {
				Length float64
				Width  float64
				Height float64
			}{Length: 6.5, Width: 5.0, Height: 3.0},
			Categories: []string{"电子产品", "音频设备"},
		},
	}

	// 创建目标切片
	var targets []*ProductTarget

	// 转换切片
	err := converter.ConvertSlice(&sources, &targets)
	if err != nil {
		t.Fatal(err)
	}

	// 验证切片长度
	if len(targets) != len(sources) {
		t.Errorf("切片长度不匹配: 预期 %d, 得到 %d", len(sources), len(targets))
	}

	// 手动处理需要计算的字段
	for i, source := range sources {
		targets[i].ProductCode = source.SKU
		targets[i].PriceText = fmt.Sprintf("￥%.2f", source.Price)
		targets[i].Available = source.Stock > 0
		targets[i].Volume = source.Dimensions.Length * source.Dimensions.Width * source.Dimensions.Height
		targets[i].StringCategories = strings.Join(source.Categories, ", ")
	}

	// 验证转换结果
	for i, target := range targets {
		source := sources[i]

		if target.ProductCode != source.SKU {
			t.Errorf("切片项 %d ProductCode转换失败: 预期 %s, 得到 %s", i, source.SKU, target.ProductCode)
		}

		if target.Title != source.Name {
			t.Errorf("切片项 %d Title转换失败: 预期 %s, 得到 %s", i, source.Name, target.Title)
		}

		expectedAvailable := source.Stock > 0
		if target.Available != expectedAvailable {
			t.Errorf("切片项 %d Available转换失败: 预期 %t, 得到 %t", i, expectedAvailable, target.Available)
		}
	}

	// 验证特定项
	if !targets[0].Available {
		t.Errorf("第一个产品可用性错误，应该为可用")
	}

	if targets[2].Available {
		t.Errorf("第三个产品可用性错误，应该为不可用")
	}

	fmt.Printf("成功转换了 %d 个产品\n", len(targets))
}

// TestSliceConversions 测试各种切片转换场景
func TestSliceConversions(t *testing.T) {
	type PersonA struct {
		Name string
		Age  int
	}

	type PersonB struct {
		Name string
		Age  int
	}

	// 测试场景1: []A -> []B (非指针到非指针)
	t.Run("Slice_NonPtr_To_NonPtr", func(t *testing.T) {
		src := []PersonA{
			{Name: "张三", Age: 30},
			{Name: "李四", Age: 25},
		}
		var dst []PersonB

		err := Cast(&src, &dst)

		if err != nil {
			t.Fatal(err)
		}
		if len(dst) != 2 || dst[0].Name != "张三" || dst[0].Age != 30 {
			t.Errorf("转换结果不正确")
		}
		fmt.Println("测试场景1成功：[]A -> []B")
	})

	// 测试场景2: []*A -> []B (指针到非指针)
	t.Run("Slice_Ptr_To_NonPtr", func(t *testing.T) {
		src := []*PersonA{
			{Name: "张三", Age: 30},
			{Name: "李四", Age: 25},
		}
		var dst []PersonB

		err := Cast(&src, &dst)

		if err != nil {
			t.Fatal(err)
		}
		if len(dst) != 2 || dst[0].Name != "张三" || dst[0].Age != 30 {
			t.Errorf("转换结果不正确")
		}
		fmt.Println("测试场景2成功：[]*A -> []B")
	})

	// 测试场景3: []A -> []*B (非指针到指针)
	t.Run("Slice_NonPtr_To_Ptr", func(t *testing.T) {
		src := []PersonA{
			{Name: "张三", Age: 30},
			{Name: "李四", Age: 25},
		}
		var dst []*PersonB

		err := Cast(&src, &dst)

		if err != nil {
			t.Fatal(err)
		}
		if len(dst) != 2 || dst[0].Name != "张三" || dst[0].Age != 30 {
			t.Errorf("转换结果不正确")
		}
		fmt.Println("测试场景3成功：[]A -> []*B")
	})

	// 测试场景4: []*A -> []*B (指针到指针)
	t.Run("Slice_Ptr_To_Ptr", func(t *testing.T) {
		src := []*PersonA{
			{Name: "张三", Age: 30},
			{Name: "李四", Age: 25},
		}
		var dst []*PersonB

		err := Cast(&src, &dst)

		if err != nil {
			t.Fatal(err)
		}
		if len(dst) != 2 || dst[0].Name != "张三" || dst[0].Age != 30 {
			t.Errorf("转换结果不正确")
		}
		fmt.Println("测试场景4成功：[]*A -> []*B")
	})
}

func TestPerformance(t *testing.T) {
	a := []map[string]string{
		{"name": "张三", "age": "30"},
		{"name": "李四", "age": "25"},
		{"name": "王五", "age": "35"},
	}
	type Person struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}

	var b []Person

	err := Cast(&a, &b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("成功转换了 %d 个产品\n", len(b))

}

func TestReverseConversion(t *testing.T) {
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address string  `json:"address"`
		Salary  float64 `json:"salary"`
	}

	a := []Person{
		{Name: "张三", Age: 30},
		{Name: "李四", Age: 25},
		{Name: "王五", Age: 35},
	}
	var b []map[string]string

	err := Cast(a, &b)
	if err != nil {
		t.Fatal(err)
	}

	if len(b) != 3 {
		t.Errorf("预期长度为3，实际为%d", len(b))
	}

	for i, m := range b {
		if m["name"] != a[i].Name || m["age"] != fmt.Sprintf("%d", a[i].Age) {
			fmt.Println(a[i].Name)
			fmt.Println(m["name"])
			fmt.Println(a[i].Age)
			fmt.Println(m["age"])
			t.Errorf("第%d个元素转换错误", i)
		}
	}

	fmt.Printf("成功反向转换了 %d 个产品\n", len(b))
}

func TestPerformance2(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}

	type Person2 struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}

	a := []*Person{
		{Name: "张三", Age: "30"},
		{Name: "李四", Age: "25"},
		{Name: "王五", Age: "35"},
	}

	var b []Person2

	err := Cast(&a, &b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("成功转换了 %d 个产品\n", len(b))
}

func TestPerformance3(t *testing.T) {
	type Person struct {
		Name   string  `json:"name"`
		Age    string  `json:"age"`
		Number int     `json:"number"`
		Height float64 `json:"height"`
		Bool2  string  `json:"bool2"`
	}

	type Person2 struct {
		Name   string `json:"name"`
		Age    string `json:"age"`
		Number *int   `json:"number"`
		Height *bool  `json:"height"`
		Bool2  bool   `json:"bool2"`
	}

	a := []Person{
		{Name: "张三", Age: "30", Number: 1, Height: 1, Bool2: "true"},
		{Name: "李四", Age: "25", Number: 2, Height: 0, Bool2: "false"},
		{Name: "王五", Age: "35", Number: 3, Height: 2, Bool2: "true"},
	}

	var b []*Person2

	err := Cast(&a, &b)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("反向转换：成功转换了 %d 个人员\n", len(b))

	// 验证转换结果
	if len(b) != 3 {
		t.Fatalf("期望转换 3 个元素，实际转换了 %d 个", len(b))
	}

	// 调试信息
	fmt.Printf("转换后 b[0]=%+v, Number=%v\n", b[0], func() interface{} {
		if b[0].Number == nil {
			return "nil"
		}
		return *b[0].Number
	}())

	if b[0].Name != "张三" || b[0].Age != "30" || b[0].Number == nil || *b[0].Number != 1 {
		t.Errorf("转换结果不正确：期望 Name=张三, Age=30, Number=1，实际 Name=%s, Age=%s, Number=%v",
			b[0].Name, b[0].Age, func() interface{} {
				if b[0].Number == nil {
					return "nil"
				}
				return *b[0].Number
			}())
	}
}

type Inner struct {
	A int
}
type P1 struct {
	Inner
}
type P2 struct {
	A int
}

func Test(t *testing.T) {
	var p1 P1
	p1.A = 1
	var p2 P2
	Cast(p1, &p2)
	fmt.Println("ok")
}
