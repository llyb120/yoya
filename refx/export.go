package refx

import (
	"fmt"

	reflect "github.com/goccy/go-reflect"
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
	refv, err := _reflectCache.getFieldByName(obj, fieldName)
	if err != nil {
		return nil, err
	}
	return refv.Interface(), nil
}

func Call(obj any, methodName string, args ...any) (result []any, finalErr error) {
	defer func() {
		if err := recover(); err != nil {
			finalErr = fmt.Errorf("call method %s failed: %v", methodName, err)
		}
	}()
	refv, err := _reflectCache.getMethodByName(obj, methodName)
	if err != nil {
		return nil, err
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
	Set func(value any) error
	Get func() (any, error)
}

type MethodHandler struct {
	Call func(args ...any) ([]any, error)
}

func GetFields(obj any) map[string]FieldHandler {
	var analysis = _reflectCache.analyze(obj)
	if analysis == nil {
		return make(map[string]FieldHandler)
	}
	var fields = make(map[string]FieldHandler)
	for name, _ := range analysis.fields {
		name := name
		fields[name] = FieldHandler{
			Set: func(value any) error {
				return Set(obj, name, value)
			},
			Get: func() (any, error) {
				return Get(obj, name)
			},
		}
	}
	return fields
}

func GetMethods(obj any) map[string]MethodHandler {
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
	return methods
}
