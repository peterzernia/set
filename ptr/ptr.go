package ptr

// Int64 turns an int64 into a *int64
func Int64(i int64) *int64 {
	p := i
	return &p
}

// Bool turns an bool into a *bool
func Bool(i bool) *bool {
	p := i
	return &p
}
