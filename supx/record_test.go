package supx

import (
	"encoding/json"
	"testing"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRecord(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  20,
	}
	record := NewRecord(person)
	record.Put("city", "New York")

	jsonBs, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonBs))

}
