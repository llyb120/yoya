package refx

import (
	"fmt"
	"testing"
)

type Inner struct {
	InnerField1 string
	InnerField2 int
}

func (i *Inner) GetInnerField1() string {
	return i.InnerField1
}

type Test struct {
	Field1 string
	Field2 int
	Test   func()
	*Inner
}

func (t *Test) GetField1() string {
	return t.Field1
}

func TestReflect(t *testing.T) {
	var test = &Test{
		Field1: "test",
		Field2: 1,
		Test: func() {
			fmt.Println("hei test")
		},
	}
	// err := Set(test, "Field1", "test2")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// err = Set(test, "Field2", "222")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if test.Field1 != "test2" {
	// 	t.Fatal("Field1 is not test2")
	// }
	// if test.Field2 != 222 {
	// 	t.Fatal("Field2 is not 1")
	// }
	// err := Set(test, "Inner.InnerField1", "inner1")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if test.Inner.InnerField1 != "inner1" {
	// 	t.Fatal("InnerField1 is not inner1")
	// }

	fields := GetFields(test, IgnoreFunc)

	// fmt.Println(fields["Field1"].Get())
	// fmt.Println(fields["Field2"].Get())
	fmt.Println(fields["InnerField1"].Get())
	// fmt.Println(fields["InnerField2"].Get())
	// fn, err := fields["Test"].Get()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fn.(func())()

	res, _ := Call(test, "GetField1")
	fmt.Println(res[0])

	methods := GetMethods(test, IncludeFieldFunc)
	fmt.Println(methods["GetInnerField1"].Call())

	methods["Test"].Call()
}
