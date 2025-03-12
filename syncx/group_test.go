package syncx

import (
	"fmt"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	var g Group
	g.Go(func() error {
		time.Sleep(1 * time.Second)
		return fmt.Errorf("error 1")
	})
	g.Go(func() error {
		time.Sleep(1 * time.Second)
		return fmt.Errorf("error 2")
	})
	err := g.Wait()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Error() != "error 1\nerror 2" {
		t.Errorf("Expected error string, got %s", err.Error())
	}
}
