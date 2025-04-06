package refx

import (
	"fmt"

	reflect "github.com/goccy/go-reflect"
	"github.com/llyb120/yoya/internal"
)

var _reflectCache = newReflectCache()

func Set(obj any, fieldName string, value any) (finalErr error) {
	defer func() {
		if err := recover(); err != nil {
			finalErr = fmt.Errorf("set field %s failed: %v", fieldName, err)
		}
	}()
	return _reflectCache.set(obj, fieldName, value)
}

func Get(obj any, fieldName string) (result any, finalErr error) {
	defer func() {
		if err := recover(); err != nil {
			finalErr = fmt.Errorf("get field %s failed: %v", fieldName, err)
		}
	}()
	refv, ok := _reflectCache.getFieldByName(obj, fieldName)
	if !ok {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}
	return refv.Interface(), nil
}

func Call(obj any, methodName string, args ...any) (result []any, finalErr error) {
	defer func() {
		if err := recover(); err != nil {
			finalErr = fmt.Errorf("call method %s failed: %v", methodName, err)
		}
	}()
	refv, ok := _reflectCache.getMethodByName(obj, methodName)
	if !ok {
		// 如果找不到，找字段
		refv, ok = _reflectCache.getFieldByName(obj, methodName)
		if !ok || refv.Kind() != reflect.Func {
			return nil, fmt.Errorf("method %s not found", methodName)
		}
	}

	var refArgs = make([]reflect.Value, len(args))
	for i, arg := range args {
		refArgs[i] = reflect.ValueOf(arg)
	}
	res := refv.Call(refArgs)
	var results = make([]any, len(res))
	for i, r := range res {
		results[i] = r.Interface()
	}
	return results, nil
}

type FieldHandler struct {
	Set  func(value any) error
	Get  func() (any, error)
	Type reflect.Type
}

type MethodHandler struct {
	Call func(args ...any) ([]any, error)
}

type fieldOption int

const (
	IgnoreFunc fieldOption = iota
	IncludeFieldFunc
)

func GetFields(obj any, opts ...fieldOption) map[string]FieldHandler {
	var analysis = _reflectCache.analyze(obj)
	if analysis == nil {
		return make(map[string]FieldHandler)
	}
	var ignoreFunc bool
	for _, opt := range opts {
		if opt == IgnoreFunc {
			ignoreFunc = true
		}
	}
	var fields = make(map[string]FieldHandler)
	for name, _ := range analysis.fields {
		name := name
		if ignoreFunc && analysis.fields[name].typ.Kind() == reflect.Func {
			continue
		}
		fields[name] = FieldHandler{
			Set: func(value any) error {
				return Set(obj, name, value)
			},
			Get: func() (any, error) {
				return Get(obj, name)
			},
			Type: analysis.fields[name].typ,
		}
	}
	return fields
}

func GetMethods(obj any, opts ...fieldOption) map[string]MethodHandler {
	var analysis = _reflectCache.analyze(obj)
	if analysis == nil {
		return make(map[string]MethodHandler)
	}
	var methods = make(map[string]MethodHandler)
	for name, _ := range analysis.methods {
		name := name
		methods[name] = MethodHandler{
			Call: func(args ...any) ([]any, error) {
				return Call(obj, name, args...)
			},
		}
	}
	var includeFunc bool
	for _, opt := range opts {
		if opt == IncludeFieldFunc {
			includeFunc = true
		}
	}
	if includeFunc {
		for name, _ := range analysis.fields {
			name := name
			if analysis.fields[name].typ.Kind() == reflect.Func {
				methods[name] = MethodHandler{
					Call: func(args ...any) ([]any, error) {
						return Call(obj, name, args...)
					},
				}
			}
		}
	}
	return methods
}

func UnsafeSetFieldValue(field reflect.Value, value reflect.Value, forceCheckType bool) {
	internal.UnsafeSetFieldValue(field, value, forceCheckType)
}
