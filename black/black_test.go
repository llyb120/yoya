package black

import (
	"fmt"
	"testing"
)

func TestToBytes(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	bs, err := ToBytes(&User{Name: "张三", Age: 18})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bs: %v \n", bs)

	user, err := FromBytes[User](bs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("user: %v \n", user)

	bs2, err := ToBytes(&User{Name: "李四", Age: 20})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bs2: %v \n", bs2)

	user2, err := FromBytes[User](bs2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("user2: %v \n", user2)

	arr := []map[string]string{{"a": "1", "b": "2"}, {"c": "3", "d": "4"}}
	bs3, err := ToBytes(&arr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bs3: %v \n", bs3)

	arr2, err := FromBytes[[]map[string]string](bs3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("arr2: %v \n", arr2)

	mp := map[string]int{"a": 1, "b": 2, "c": 3}
	bs4, err := ToBytes(mp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("bs4: %v \n", bs4)

	mp2, err := FromBytes[map[string]int](bs4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("mp2: %v \n", mp2)

}
