package supx

import (
	"fmt"
	"testing"

	_ "github.com/llyb120/yoya/objx"
)

func TestData(t *testing.T) {
	d := NewData[int]()
	d.Set(1)
	fmt.Println(d.Clone())
}
