package ptr

func String(value string) *string {
	return &value
}

func Optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func Deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func Int64(value int64) *int64 {
	return &value
}

func Float64(value float64) *float64 {
	return &value
}
