package internal

import "time"

type IGroup interface {
	Go(func() error)
	Wait(timeout ...time.Duration) error
}
