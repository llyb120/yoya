package objx

import (
	"fmt"

	reflect "github.com/goccy/go-reflect"
	"github.com/llyb120/yoya/errx"
	"github.com/llyb120/yoya/internal"
)

func Assign[T comparable, K0 any, K1 any](dest map[T]K0, source map[T]K1) {
	var zero0 K0
	var zero1 K1
	isSameType := reflect.TypeOf(zero0) == reflect.TypeOf(zero1)
	if isSameType {
		for k, v := range source {
			dest[k] = any(v).(K0)
		}
		return
	} else {
		// 一次性转过来
		var mp map[T]K0
		if err := internal.Cast(&mp, source); err != nil {
			fmt.Println("err:", err)
			return
		}
		for k, v := range mp {
			dest[k] = v
		}
		return
	}
}

// 确保某些变量是某种值
// 参数必须是偶数，第n个参数必须是指针，n+1参数是原始的值
// 如果转换成功，则会返回true
// 如果任意一个转换失败，则会返回false
// 例如 a = 1 b = 2 c = 3 a/b/c均为any类型
// var d,e,f int
// Ensure(&d, a, &e, b, &f, c)
func Ensure(objs ...any) bool {
	// 如果不是偶数，pass
	if len(objs)%2 != 0 {
		return false
	}

	for i := 0; i < len(objs); i += 2 {
		if !ensure(objs[i], objs[i+1]) {
			return false
		}
	}
	return true
}

func ensure(obj, target any) bool {
	err := errx.Try(func() error {
		// 第一个必须是一个指针
		vf := reflect.ValueOf(obj)
		if vf.Kind() != reflect.Ptr {
			return fmt.Errorf("not paired")
		}
		targetVf := reflect.ValueOf(target)
		// 转换成目标类型
		for targetVf.Kind() == reflect.Ptr || targetVf.Kind() == reflect.Interface {
			targetVf = targetVf.Elem()
		}
		targetElem := targetVf.Convert(vf.Elem().Type())
		vf.Elem().Set(targetElem)
		return nil
	})
	return err == nil
}
