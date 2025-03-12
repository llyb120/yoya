package errx

import (
	"strings"
	"sync"
)

type MultiError struct {
	mu   sync.RWMutex
	errs []error
}

func (e *MultiError) Error() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var errs []string
	for _, err := range e.errs {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, "\n")
}

func (e *MultiError) Add(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.errs = append(e.errs, err)
}

func (e *MultiError) HasError() bool {
	return len(e.errs) > 0
}
