package objx

func Or[T any](v T, def T) T {
	if any(v) == nil {
		return def
	}
	// 如果是string
	if s, ok := any(v).(string); ok && s == "" {
		return def
	}
	return v
}
