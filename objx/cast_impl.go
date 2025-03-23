package objx

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	reflect "github.com/goccy/go-reflect"
)

// 缓存相关的结构体定义
type typeCache struct {
	reflectType reflect.Type
	fields      map[string]*fieldCache
}

type fieldCache struct {
	field     reflect.StructField
	jsonTag   string
	fieldType reflect.Type
	indexes   []int // 字段索引，支持嵌套结构体
}

// 字段映射缓存，用于缓存源类型到目标类型的字段映射关系
type mappingCacheItem struct {
	srcFieldCache  *fieldCache
	destFieldCache *fieldCache
}
type mappingCache = []*mappingCacheItem

// Converter 是一个基于 reflect 的结构体转换器
type Converter struct {
	// 缓存各类型的反射信息，提高性能
	typeCache map[reflect.Type]*typeCache
	// 缓存类型映射关系
	mappingCache map[reflect.Type]map[reflect.Type]mappingCache
	// 读写锁，保证并发安全
	rwMutex sync.RWMutex
}

// newConverter 创建一个新的转换器实例
func newConverter() *Converter {
	return &Converter{
		typeCache:    make(map[reflect.Type]*typeCache),
		mappingCache: make(map[reflect.Type]map[reflect.Type]mappingCache),
	}
}

// 获取或创建类型缓存
func (c *Converter) getOrCreateTypeCache(typ reflect.Type) *typeCache {
	// 先尝试用读锁获取缓存
	c.rwMutex.RLock()
	cached, ok := c.typeCache[typ]
	c.rwMutex.RUnlock()

	if ok {
		return cached
	}

	// 缓存不存在，需要创建新的缓存，使用写锁
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 双重检查，防止其他协程已经创建了缓存
	if cached, ok := c.typeCache[typ]; ok {
		return cached
	}

	// 创建新的类型缓存
	cache := &typeCache{
		reflectType: typ,
		fields:      make(map[string]*fieldCache),
	}

	// 如果是结构体类型，预缓存所有字段
	if typ.Kind() == reflect.Struct {
		structType := typ
		// 缓存普通字段
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldName := field.Name

			// 获取字段的JSON tag
			jsonTag := ""
			tag := field.Tag.Get("json")
			if tag != "" {
				// 解析json tag，去除omitempty等选项
				parts := strings.Split(tag, ",")
				if parts[0] != "-" { // 忽略 json:"-" 的字段
					jsonTag = parts[0]
				}
			}

			// 建立字段缓存
			cache.fields[fieldName] = &fieldCache{
				field:     field,
				jsonTag:   jsonTag,
				fieldType: field.Type,
				indexes:   []int{i}, // 记录字段索引
			}

			// 处理嵌套结构体字段
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				c.processEmbeddedStruct(field.Type, cache, []int{i})
			}
		}
	}

	c.typeCache[typ] = cache
	return cache
}

// 处理嵌套的结构体字段
func (c *Converter) processEmbeddedStruct(embedType reflect.Type, parentCache *typeCache, parentIndexes []int) {
	for i := 0; i < embedType.NumField(); i++ {
		field := embedType.Field(i)
		fieldName := field.Name

		// 跳过非导出字段
		if field.PkgPath != "" {
			continue
		}

		// 获取字段的JSON tag
		jsonTag := ""
		tag := field.Tag.Get("json")
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "-" {
				jsonTag = parts[0]
			}
		}

		// 创建完整的索引路径
		indexes := make([]int, len(parentIndexes)+1)
		copy(indexes, parentIndexes)
		indexes[len(parentIndexes)] = i

		// 检查是否已存在同名字段
		if _, exists := parentCache.fields[fieldName]; !exists {
			parentCache.fields[fieldName] = &fieldCache{
				field:     field,
				jsonTag:   jsonTag,
				fieldType: field.Type,
				indexes:   indexes,
			}
		}

		// 递归处理嵌套结构体
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			c.processEmbeddedStruct(field.Type, parentCache, indexes)
		}
	}
}

// 获取类型之间的字段映射关系
func (c *Converter) getOrCreateMappingCache(srcType, destType reflect.Type, srcCache, destCache *typeCache) mappingCache {
	// 先尝试用读锁获取缓存
	c.rwMutex.RLock()
	cached0, ok := c.mappingCache[srcType]
	var cached mappingCache
	if ok {
		cached = cached0[destType]
	}
	c.rwMutex.RUnlock()

	if cached != nil {
		return cached
	}

	// 缓存不存在，需要创建新的缓存，使用写锁
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	// 双重检查，防止其他协程已经创建了缓存
	if cached0, ok := c.mappingCache[srcType]; ok {
		if cached = cached0[destType]; cached != nil {
			return cached
		}
	}

	// 创建新的映射缓存
	var cache mappingCache
	// cache := make(mappingCache)

	// // 确保srcType对应的map存在
	if _, ok := c.mappingCache[srcType]; !ok {
		c.mappingCache[srcType] = make(map[reflect.Type]mappingCache)
	}

	for i := 0; i < destType.NumField(); i++ {
		destField := destType.Field(i)
		destFieldName := destField.Name

		// 1. 首先检查目标字段是否有json tag
		destFieldCache := destCache.fields[destFieldName]

		destJsonTag := destFieldCache.jsonTag

		var srcFieldCacheFound *fieldCache
		var srcFieldFoundExists bool

		if destJsonTag != "" {
			// 如果目标字段有json tag，尝试在源结构体中查找具有相同json tag的字段
			for _, srcCache := range srcCache.fields {
				if srcCache.jsonTag == destJsonTag {
					// 找到了具有相同json tag的源字段
					srcFieldCacheFound = srcCache
					srcFieldFoundExists = true
					break
				}
			}

			// 如果没有找到具有相同json tag的字段，尝试查找具有相同名称的字段
			if !srcFieldFoundExists {
				for _, srcCache := range srcCache.fields {
					if srcCache.jsonTag == destFieldName {
						// 源字段名与目标字段的json tag匹配
						srcFieldCacheFound = srcCache
						srcFieldFoundExists = true
						break
					}
				}
			}
		}

		// 2. 如果通过json tag没有找到匹配的字段，尝试直接通过字段名匹配
		if !srcFieldFoundExists {
			if srcCache, found := srcCache.fields[destFieldName]; found {
				srcFieldCacheFound = srcCache
				srcFieldFoundExists = true
			}
		}

		// 如果没有找到匹配，跳过此字段
		if !srcFieldFoundExists {
			continue
		}

		// 缓存映射关系
		cache = append(cache, &mappingCacheItem{
			srcFieldCache:  srcFieldCacheFound,
			destFieldCache: destFieldCache,
		})
	}

	c.mappingCache[srcType][destType] = cache
	return cache
}

// Convert 将源类型转换为目标类型
func (c *Converter) Convert(src, dest interface{}) error {
	// 空值检查
	if dest == nil {
		return errors.New("源或目标不能为空")
	}

	// 获取源值和目标值
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// 检查目标是否为指针
	if destValue.Kind() != reflect.Ptr {
		return errors.New("目标必须是指针类型")
	}

	// 获取源类型和目标类型
	srcType := reflect.TypeOf(src)
	destType := reflect.TypeOf(dest)

	// 处理源类型
	var srcElemTypeObj reflect.Type
	isSrcPtr := srcValue.Kind() == reflect.Ptr
	if isSrcPtr {
		// 如果是指针，获取其指向的类型
		srcElemTypeObj = srcType.Elem()
		// 解引用源值
		srcValue = srcValue.Elem()
	} else {
		srcElemTypeObj = srcType
	}

	// 处理目标类型（目标已经确认是指针）
	destElemTypeObj := destType.Elem()
	destElemValue := destValue.Elem()

	// 特殊处理： 如果源是 map 类型，目标是结构体类型
	if srcValue.Kind() == reflect.Map && destElemValue.Kind() == reflect.Struct {
		return c.convertMapToStruct(srcValue, destElemValue, srcElemTypeObj, destElemTypeObj)
	}

	// 特殊处理： 如果源是结构体类型，目标是 map 类型
	if srcValue.Kind() == reflect.Struct && destElemValue.Kind() == reflect.Map {
		return c.convertStructToMap(srcValue, destElemValue, srcElemTypeObj, destElemTypeObj)
	}

	// 如果源和目标不是结构体，尝试直接转换值
	if srcValue.Kind() != reflect.Struct || destElemValue.Kind() != reflect.Struct {
		return c.convertNonStructValues(srcValue, destElemValue, srcElemTypeObj, destElemTypeObj)
	}

	// 确保它们是结构体类型
	if srcElemTypeObj.Kind() != reflect.Struct || destElemTypeObj.Kind() != reflect.Struct {
		return errors.New("源和目标必须是结构体或结构体指针")
	}

	// 从缓存获取类型信息
	srcCache := c.getOrCreateTypeCache(srcElemTypeObj)
	destCache := c.getOrCreateTypeCache(destElemTypeObj)
	_ = destCache // 避免未使用警告

	// 获取类型映射缓存
	mappingCache := c.getOrCreateMappingCache(srcElemTypeObj, destElemTypeObj, srcCache, destCache)

	// 获取目标结构体的类型
	// destStructType := destElemTypeObj

	// 遍历目标结构体的所有字段
	for _, mappingCacheItem := range mappingCache {

		// 获取源字段值 - 使用字段索引而不是名称
		srcFieldReflect := srcValue.FieldByIndex(mappingCacheItem.srcFieldCache.indexes)

		if !srcFieldReflect.IsValid() {
			continue
		}

		srcFieldValue := srcFieldReflect.Interface()
		// 获取目标字段值 - 也使用索引访问
		destFieldReflect := destElemValue.FieldByIndex(mappingCacheItem.destFieldCache.indexes)

		// 获取字段类型
		destFieldType := mappingCacheItem.destFieldCache.fieldType
		srcFieldType := mappingCacheItem.srcFieldCache.fieldType

		// 如果源字段和目标字段类型相同，直接赋值
		if srcFieldType.String() == destFieldType.String() {
			c.unsafeSetField(destElemValue, mappingCacheItem.destFieldCache.indexes, srcFieldReflect)
			continue
		}

		// 基本类型转换规则
		if c.canConvert(srcFieldType, destFieldType) {
			convertedValue, err := c.convertValue(srcFieldValue, srcFieldType, destFieldType)
			if err == nil && convertedValue != nil {
				convertedReflect := reflect.ValueOf(convertedValue)

				// 处理类型不匹配的情况
				if !convertedReflect.Type().AssignableTo(destFieldReflect.Type()) {
					// 如果目标是指针，但转换后的值不是指针
					if destFieldReflect.Kind() == reflect.Ptr && convertedReflect.Kind() != reflect.Ptr {
						// 创建一个新的指针
						ptrValue := reflect.New(convertedReflect.Type())
						// 设置指针指向的值
						ptrValue.Elem().Set(convertedReflect)
						destFieldReflect.Set(ptrValue)
					} else if destFieldReflect.Kind() != reflect.Ptr && convertedReflect.Kind() == reflect.Ptr {
						// 如果目标不是指针，但转换后的值是指针
						destFieldReflect.Set(convertedReflect.Elem())
					} else if convertedReflect.Type().ConvertibleTo(destFieldReflect.Type()) {
						// 如果类型不匹配但可以转换
						destFieldReflect.Set(convertedReflect.Convert(destFieldReflect.Type()))
					}
				} else {
					// 类型匹配，直接赋值
					c.unsafeSetField(destElemValue, mappingCacheItem.destFieldCache.indexes, convertedReflect)
				}
			}
		}
	}

	return nil
}

// 使用unsafe设置字段值，性能更高但需谨慎使用
func (c *Converter) unsafeSetField(structValue reflect.Value, fieldIndex []int, value reflect.Value) {
	// 获取目标字段
	field := structValue.FieldByIndex(fieldIndex)
	c.unsafeSetFieldValue(field, value)
}

func (c *Converter) unsafeSetFieldValue(field reflect.Value, value reflect.Value) {
	// 如果目标字段和值的类型相同，则使用unsafe直接设置
	if field.Type() == value.Type() && field.CanAddr() && value.CanAddr() {
		// 获取目标字段的指针
		fieldPtr := unsafe.Pointer(field.UnsafeAddr())

		// 获取值的指针
		valuePtr := unsafe.Pointer(value.UnsafeAddr())

		// 计算字段大小
		size := field.Type().Size()

		// 直接复制内存
		if size > 0 {
			c.typedmemmove(field.Type(), fieldPtr, valuePtr)
		}
	} else {
		// 类型不同，回退到使用标准反射
		field.Set(value)
	}
}

// 通过unsafe实现类似reflect.typedmemmove的功能
// 这个函数假设src和dst都有效，且大小匹配
func (c *Converter) typedmemmove(typ reflect.Type, dst, src unsafe.Pointer) {
	// 使用非导出的内存复制
	// 实际实现中，我们应该参考runtime.typedmemmove的实现
	// 这里为了简化，直接使用标准库的copy
	size := typ.Size()
	c.memmove(dst, src, size)
}

// 封装底层内存复制操作
func (c *Converter) memmove(dst, src unsafe.Pointer, size uintptr) {
	// 转为切片进行复制
	dstSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(dst),
		Len:  int(size),
		Cap:  int(size),
	}))

	srcSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(src),
		Len:  int(size),
		Cap:  int(size),
	}))

	copy(dstSlice, srcSlice)
}

// convertMapToStruct 将 map 转换为结构体
func (c *Converter) convertMapToStruct(srcMap reflect.Value, destStruct reflect.Value, srcType, destType reflect.Type) error {
	// 检查 map 的键类型是否为 string
	if srcMap.Type().Key().Kind() != reflect.String {
		return errors.New("map 的键类型必须是 string")
	}

	// 从缓存中获取目标类型信息，如果没有则创建
	destCache := c.getOrCreateTypeCache(destType)

	// 遍历结构体的每个字段
	for fieldName, destFieldCache := range destCache.fields {
		// 获取目标字段信息，用于调试和日志
		jsonTag := destFieldCache.jsonTag

		// 确定要从 map 中获取的键名
		mapKey := fieldName
		if jsonTag != "" {
			mapKey = jsonTag // 优先使用 json tag
		}

		// 从 map 中获取值
		mapValue := srcMap.MapIndex(reflect.ValueOf(mapKey))
		if !mapValue.IsValid() {
			// 如果 map 中没有对应的键，尝试使用字段名作为键
			if jsonTag != "" {
				mapValue = srcMap.MapIndex(reflect.ValueOf(fieldName))
			}

			// 如果仍然找不到，跳过此字段
			if !mapValue.IsValid() {
				continue
			}
		}

		// 获取目标字段
		destFieldValue := destStruct.FieldByIndex(destFieldCache.indexes)
		if !destFieldValue.IsValid() || !destFieldValue.CanSet() {
			continue
		}

		// 处理不同类型的值
		if mapValue.Type().AssignableTo(destFieldValue.Type()) {
			// 类型匹配，直接赋值
			c.unsafeSetField(destStruct, destFieldCache.indexes, mapValue)
		} else {
			// 类型不匹配，尝试转换
			convertedValue, err := c.convertValue(mapValue.Interface(), reflect.TypeOf(mapValue.Interface()), reflect.TypeOf(destFieldValue.Interface()))
			if err == nil && convertedValue != nil {
				convertedReflect := reflect.ValueOf(convertedValue)
				c.unsafeSetField(destStruct, destFieldCache.indexes, convertedReflect)
			}
		}
	}

	return nil
}

// convertStructToMap 将结构体转换为 map
func (c *Converter) convertStructToMap(srcStruct reflect.Value, destMap reflect.Value, srcType, destType reflect.Type) error {
	// 检查 map 的键类型是否为 string
	if destType.Key().Kind() != reflect.String {
		return errors.New("map 的键类型必须是 string")
	}

	// 如果 map 是 nil，初始化它
	if destMap.IsNil() {
		destMap.Set(reflect.MakeMap(destType))
	}

	// 从缓存中获取源类型信息，如果没有则创建
	srcCache := c.getOrCreateTypeCache(srcType)

	// 遍历结构体的每个字段
	for fieldName, srcFieldCache := range srcCache.fields {
		// 获取字段值
		srcFieldValue := srcStruct.FieldByIndex(srcFieldCache.indexes)
		if !srcFieldValue.IsValid() {
			continue
		}

		// 确定要存入 map 的键名
		mapKey := fieldName
		if srcFieldCache.jsonTag != "" {
			mapKey = srcFieldCache.jsonTag // 优先使用 json tag
		}

		// 存入 map
		destMapElemType := destType.Elem()

		// 如果字段值可以直接赋值给 map 值类型
		if srcFieldValue.Type().AssignableTo(destMapElemType) {
			destMap.SetMapIndex(reflect.ValueOf(mapKey), srcFieldValue)
		} else {
			// 类型不匹配，尝试转换
			convertedValue, err := c.convertValue(srcFieldValue.Interface(), reflect.TypeOf(srcFieldValue.Interface()), reflect.TypeOf(reflect.Zero(destMapElemType).Interface()))
			if err == nil && convertedValue != nil {
				destMap.SetMapIndex(reflect.ValueOf(mapKey), reflect.ValueOf(convertedValue))
			} else {
				// 如果无法转换，尝试使用 fmt.Sprint
				stringValue := fmt.Sprintf("%v", srcFieldValue.Interface())
				if destMapElemType.Kind() == reflect.String {
					destMap.SetMapIndex(reflect.ValueOf(mapKey), reflect.ValueOf(stringValue))
				}
			}
		}
	}

	return nil
}

// convertNonStructValues 处理非结构体之间的转换
func (c *Converter) convertNonStructValues(src, dest reflect.Value, srcType, destType reflect.Type) error {
	// 如果源和目标类型相同，直接赋值
	if srcType.AssignableTo(destType) {
		if dest.CanSet() {
			c.unsafeSetField(dest, []int{0}, src)
		}
		return nil
	}

	// 尝试进行类型转换
	if dest.CanSet() {
		convertedValue, err := c.convertValue(src.Interface(), srcType, destType)
		if err == nil && convertedValue != nil {
			// 处理转换后的值可能与目标类型不匹配的情况
			convertedValueReflect := reflect.ValueOf(convertedValue)

			// 如果目标是指针，但转换后的值不是指针
			if dest.Kind() == reflect.Ptr && convertedValueReflect.Kind() != reflect.Ptr {
				// 创建一个新的指针
				ptrValue := reflect.New(convertedValueReflect.Type())
				// 设置指针指向的值
				ptrValue.Elem().Set(convertedValueReflect)
				c.unsafeSetField(dest, []int{0}, ptrValue)
				return nil
			}

			// 如果目标不是指针，但转换后的值是指针
			if dest.Kind() != reflect.Ptr && convertedValueReflect.Kind() == reflect.Ptr {
				// 取指针指向的值
				c.unsafeSetFieldValue(dest, convertedValueReflect.Elem())
				return nil
			}

			// 如果类型可以直接赋值
			if convertedValueReflect.Type().AssignableTo(destType) {
				c.unsafeSetFieldValue(dest, convertedValueReflect)
				return nil
			}

			// 如果类型不匹配但可以转换
			if convertedValueReflect.Type().ConvertibleTo(destType) {
				c.unsafeSetFieldValue(dest, convertedValueReflect.Convert(destType))
				return nil
			}
		}
	}

	return fmt.Errorf("无法将类型 %s 转换为 %s", src.Type(), dest.Type())
}

// ConvertSlice 将源切片转换为目标切片
func (c *Converter) ConvertSlice(srcSlice, destSlice interface{}) error {
	srcValue := reflect.ValueOf(srcSlice)
	destValue := reflect.ValueOf(destSlice)

	// 检查源和目标是否为指针
	if destValue.Kind() != reflect.Ptr {
		return errors.New("源和目标切片必须是指针")
	}

	// 解引用获取切片值
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	destValue = destValue.Elem()

	// 检查源是否为切片
	if srcValue.Kind() != reflect.Slice {
		return errors.New("源必须是切片指针")
	}

	// 检查目标是否为切片
	if destValue.Kind() != reflect.Slice {
		return errors.New("目标必须是切片指针")
	}

	// 获取源切片长度
	srcLen := srcValue.Len()

	// 获取源和目标元素类型
	srcElemType := srcValue.Type().Elem()
	destElemType := destValue.Type().Elem()

	// 判断是否为结构体到 map 或 map 到结构体的转换
	isMapToStruct := srcElemType.Kind() == reflect.Map &&
		(destElemType.Kind() == reflect.Struct ||
			(destElemType.Kind() == reflect.Ptr && destElemType.Elem().Kind() == reflect.Struct))

	isStructToMap := (srcElemType.Kind() == reflect.Struct ||
		(srcElemType.Kind() == reflect.Ptr && srcElemType.Elem().Kind() == reflect.Struct)) &&
		destElemType.Kind() == reflect.Map

	// 创建新的目标切片
	newSlice := reflect.MakeSlice(destValue.Type(), srcLen, srcLen)

	// 是否是指针元素类型
	isDestElemPtr := destElemType.Kind() == reflect.Ptr
	destElemValueType := destElemType
	if isDestElemPtr {
		destElemValueType = destElemType.Elem()
	}

	// 遍历并转换每个元素
	for i := 0; i < srcLen; i++ {
		srcElem := srcValue.Index(i)

		// 创建目标元素
		var destElem reflect.Value
		if isDestElemPtr {
			// 如果目标元素是指针类型，创建一个新的指针
			destElem = reflect.New(destElemValueType)
		} else {
			// 如果目标元素不是指针类型，创建一个零值
			destElem = reflect.New(destElemType)
		}

		// 根据类型选择不同的转换方法
		var err error
		if isMapToStruct {
			// 从 map 到结构体的转换
			mapValue := srcElem
			if srcElem.Kind() == reflect.Ptr {
				mapValue = srcElem.Elem()
			}

			// 目标结构体始终是指针的 Elem
			structValue := destElem.Elem()
			if isDestElemPtr {
				// 无需额外操作，destElem 已经是指针
			} else {
				// 无需额外操作，structValue 已经是非指针
			}

			err = c.convertMapToStruct(mapValue, structValue, mapValue.Type(), structValue.Type())
		} else if isStructToMap {
			// 从结构体到 map 的转换
			structValue := srcElem
			if srcElem.Kind() == reflect.Ptr {
				structValue = srcElem.Elem()
			}

			mapValue := destElem.Elem()
			err = c.convertStructToMap(structValue, mapValue, structValue.Type(), mapValue.Type())
		} else {
			// 常规转换 - 为每个元素创建正确的指针
			var srcPtr interface{}

			// 为源元素创建指针（如果需要）
			if srcElem.Kind() != reflect.Ptr {
				// 创建一个临时变量来保存元素值
				temp := reflect.New(srcElem.Type())
				temp.Elem().Set(srcElem)
				srcPtr = temp.Interface()
			} else {
				srcPtr = srcElem.Interface()
			}

			// 调用 Convert 进行转换
			err = c.Convert(srcPtr, destElem.Interface())
		}

		if err != nil {
			return fmt.Errorf("转换切片元素 %d 失败: %s", i, err)
		}

		// 将转换后的元素设置到目标切片中
		if isDestElemPtr {
			// 如果目标元素类型是指针，直接设置
			newSlice.Index(i).Set(destElem)
		} else {
			// 如果目标元素类型不是指针，设置 Elem 的值
			newSlice.Index(i).Set(destElem.Elem())
		}
	}

	// 设置新切片到目标切片
	c.unsafeSetFieldValue(destValue, newSlice)

	return nil
}

// canConvert 判断是否可以进行类型转换
func (c *Converter) canConvert(srcType, destType reflect.Type) bool {
	// 处理指针类型
	if destType.Kind() == reflect.Ptr {
		// 获取指针指向的类型
		elemType := destType.Elem()
		// 递归检查源类型是否可以转换为指针指向的类型
		return c.canConvert(srcType, elemType)
	}

	// 如果源是指针，检查指针指向的类型
	if srcType.Kind() == reflect.Ptr {
		elemType := srcType.Elem()
		return c.canConvert(elemType, destType)
	}

	// 数值类型之间可以互相转换
	if isNumeric(srcType.Kind()) && isNumeric(destType.Kind()) {
		return true
	}

	// 字符串和数值类型可以互相转换
	if (srcType.Kind() == reflect.String && isNumeric(destType.Kind())) ||
		(isNumeric(srcType.Kind()) && destType.Kind() == reflect.String) {
		return true
	}

	// 布尔类型转换规则
	if srcType.Kind() == reflect.Bool && destType.Kind() == reflect.Bool {
		return true
	}

	// 数值类型到布尔类型的转换
	if isNumeric(srcType.Kind()) && destType.Kind() == reflect.Bool {
		return true
	}

	// 布尔类型到数值类型的转换
	if srcType.Kind() == reflect.Bool && isNumeric(destType.Kind()) {
		return true
	}

	// 字符串到布尔类型的转换
	if srcType.Kind() == reflect.String && destType.Kind() == reflect.Bool {
		return true
	}

	// 布尔类型到字符串的转换
	if srcType.Kind() == reflect.Bool && destType.Kind() == reflect.String {
		return true
	}

	// 时间类型转换规则
	srcTypeName := srcType.String()
	destTypeName := destType.String()
	if (srcTypeName == "time.Time" && destType.Kind() == reflect.String) ||
		(srcType.Kind() == reflect.String && destTypeName == "time.Time") {
		return true
	}

	return false
}

// 判断是否为数值类型
func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// convertValue 将值从一种类型转换为另一种类型
func (c *Converter) convertValue(value interface{}, srcType, destType reflect.Type) (interface{}, error) {
	// 检查空值
	if value == nil {
		return nil, nil
	}

	// 获取目标类型的 Kind
	destKind := destType.Kind()

	// 特殊处理：目标类型是指针
	if destKind == reflect.Ptr {
		// 获取指针指向的类型
		elemType := destType.Elem()
		// 递归调用 convertValue 转换为元素类型
		elemValue, err := c.convertValue(value, srcType, elemType)
		if err != nil {
			return nil, err
		}
		if elemValue == nil {
			return nil, nil
		}

		// 创建新的指针并设置值
		ptrValue := reflect.New(reflect.TypeOf(elemValue))
		ptrValue.Elem().Set(reflect.ValueOf(elemValue))
		return ptrValue.Interface(), nil
	}

	// 常规类型转换
	srcKind := srcType.Kind()

	// 如果源是指针，获取它指向的值
	if srcKind == reflect.Ptr {
		v := reflect.ValueOf(value)
		if v.IsNil() {
			return nil, nil
		}
		elem := v.Elem().Interface()
		return c.convertValue(elem, reflect.TypeOf(elem), destType)
	}

	// 使用原生reflect进行值转换
	srcReflectVal := reflect.ValueOf(value)

	// 字符串类型转换
	if destKind == reflect.String {
		// 任何类型都可以转换为字符串
		return fmt.Sprintf("%v", value), nil
	}

	// 数值类型转换
	if isNumeric(destKind) {
		var val interface{}
		var err error

		switch srcKind {
		case reflect.String:
			// 字符串转数值
			strVal := srcReflectVal.String()
			switch destKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				var intVal int64
				_, err = fmt.Sscanf(strVal, "%d", &intVal)
				if err != nil {
					return nil, err
				}
				val = intVal
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				var uintVal uint64
				_, err = fmt.Sscanf(strVal, "%d", &uintVal)
				if err != nil {
					return nil, err
				}
				val = uintVal
			case reflect.Float32, reflect.Float64:
				var floatVal float64
				_, err = fmt.Sscanf(strVal, "%f", &floatVal)
				if err != nil {
					return nil, err
				}
				val = floatVal
			}
		default:
			// 数值类型之间的转换
			switch srcReflectVal.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				val = srcReflectVal.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val = srcReflectVal.Uint()
			case reflect.Float32, reflect.Float64:
				val = srcReflectVal.Float()
			default:
				return nil, errors.New("无法转换为数值类型")
			}
		}

		// 处理具体的数值类型
		switch destKind {
		case reflect.Int:
			if i, ok := val.(int64); ok {
				return int(i), nil
			}
			if u, ok := val.(uint64); ok {
				return int(u), nil
			}
			if f, ok := val.(float64); ok {
				return int(f), nil
			}
		case reflect.Int8:
			if i, ok := val.(int64); ok {
				return int8(i), nil
			}
			if u, ok := val.(uint64); ok {
				return int8(u), nil
			}
			if f, ok := val.(float64); ok {
				return int8(f), nil
			}
		case reflect.Int16:
			if i, ok := val.(int64); ok {
				return int16(i), nil
			}
			if u, ok := val.(uint64); ok {
				return int16(u), nil
			}
			if f, ok := val.(float64); ok {
				return int16(f), nil
			}
		case reflect.Int32:
			if i, ok := val.(int64); ok {
				return int32(i), nil
			}
			if u, ok := val.(uint64); ok {
				return int32(u), nil
			}
			if f, ok := val.(float64); ok {
				return int32(f), nil
			}
		case reflect.Int64:
			if i, ok := val.(int64); ok {
				return i, nil
			}
			if u, ok := val.(uint64); ok {
				return int64(u), nil
			}
			if f, ok := val.(float64); ok {
				return int64(f), nil
			}
		case reflect.Uint:
			if i, ok := val.(int64); ok && i >= 0 {
				return uint(i), nil
			}
			if u, ok := val.(uint64); ok {
				return uint(u), nil
			}
			if f, ok := val.(float64); ok && f >= 0 {
				return uint(f), nil
			}
		case reflect.Uint8:
			if i, ok := val.(int64); ok && i >= 0 {
				return uint8(i), nil
			}
			if u, ok := val.(uint64); ok {
				return uint8(u), nil
			}
			if f, ok := val.(float64); ok && f >= 0 {
				return uint8(f), nil
			}
		case reflect.Uint16:
			if i, ok := val.(int64); ok && i >= 0 {
				return uint16(i), nil
			}
			if u, ok := val.(uint64); ok {
				return uint16(u), nil
			}
			if f, ok := val.(float64); ok && f >= 0 {
				return uint16(f), nil
			}
		case reflect.Uint32:
			if i, ok := val.(int64); ok && i >= 0 {
				return uint32(i), nil
			}
			if u, ok := val.(uint64); ok {
				return uint32(u), nil
			}
			if f, ok := val.(float64); ok && f >= 0 {
				return uint32(f), nil
			}
		case reflect.Uint64:
			if i, ok := val.(int64); ok && i >= 0 {
				return uint64(i), nil
			}
			if u, ok := val.(uint64); ok {
				return u, nil
			}
			if f, ok := val.(float64); ok && f >= 0 {
				return uint64(f), nil
			}
		case reflect.Float32:
			if i, ok := val.(int64); ok {
				return float32(i), nil
			}
			if u, ok := val.(uint64); ok {
				return float32(u), nil
			}
			if f, ok := val.(float64); ok {
				return float32(f), nil
			}
		case reflect.Float64:
			if i, ok := val.(int64); ok {
				return float64(i), nil
			}
			if u, ok := val.(uint64); ok {
				return float64(u), nil
			}
			if f, ok := val.(float64); ok {
				return f, nil
			}
		}
	}

	// 布尔类型转换
	if destKind == reflect.Bool {
		// 从源值中获取实际值
		sourceValue := value

		// 从布尔值转换
		if b, ok := sourceValue.(bool); ok {
			return b, nil
		}
		// 从整数转换 (0 为 false, 非 0 为 true)
		if i, ok := sourceValue.(int64); ok {
			return i != 0, nil
		} else if i, ok := sourceValue.(int); ok {
			return i != 0, nil
		} else if i, ok := sourceValue.(int32); ok {
			return i != 0, nil
		}
		// 从无符号整数转换
		if u, ok := sourceValue.(uint64); ok {
			return u != 0, nil
		} else if u, ok := sourceValue.(uint); ok {
			return u != 0, nil
		} else if u, ok := sourceValue.(uint32); ok {
			return u != 0, nil
		}
		// 从浮点数转换 (0.0 为 false, 非 0.0 为 true)
		if f, ok := sourceValue.(float64); ok {
			return f != 0.0, nil
		} else if f, ok := sourceValue.(float32); ok {
			return f != 0.0, nil
		}
		// 从字符串转换
		if s, ok := sourceValue.(string); ok {
			s = strings.ToLower(s)
			return s == "true" || s == "yes" || s == "1" || s == "t" || s == "y", nil
		}
	}

	return nil, errors.New("不支持的类型转换")
}
