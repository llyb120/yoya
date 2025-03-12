package errx

import (
	"errors"
	"testing"
)

func TestMultiError(t *testing.T) {
	var merr MultiError
	merr.Add(errors.New("error 1"))
	merr.Add(errors.New("error 2"))
	merr.Add(errors.New("error 3"))
	if !merr.HasError() {
		t.Errorf("Expected error, got nil")
	}
	if merr.Error() != "error 1\nerror 2\nerror 3" {
		t.Errorf("Expected error string, got %s", merr.Error())
	}
}
