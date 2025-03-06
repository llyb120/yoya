package stlx

import (
	"encoding/json"
	"testing"
)

func TestOrderedMap(t *testing.T) {
	om := NewMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	jsonData, err := json.Marshal(om)
	if err != nil {
		t.Errorf("MarshalJSON failed: %v", err)
	}
	t.Logf("jsonData: %s", jsonData)

	om.For(func(key string, value int) bool {
		t.Logf("key: %s, value: %d", key, value)
		return true
	})

	// 反序列化
	om2 := NewMap[string, int]()
	err = json.Unmarshal(jsonData, om2)
	if err != nil {
		t.Errorf("UnmarshalJSON failed: %v", err)
	}
	om2.For(func(key string, value int) bool {
		t.Logf("key: %s, value: %d", key, value)
		return true
	})

}
