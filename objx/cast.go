package objx

import "github.com/llyb120/yoya/internal"

func Cast(dest any, src any) error {
	return internal.Cast(dest, src)
}
