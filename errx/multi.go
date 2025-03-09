package errx

import "strings"

type MultiError []error

func (e *MultiError) Error() string {
	var errs []string
	for _, err := range *e {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, "\n")
}

func (e *MultiError) Add(err error) {
	*e = append(*e, err)
}

func (e *MultiError) HasError() bool {
	return len(*e) > 0
}
