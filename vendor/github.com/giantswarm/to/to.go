package to

func Int64(i *int64) int64 {
	return *i
}

func Int64P(i int64) *int64 {
	return &i
}

func String(s *string) string {
	return *s
}

func StringP(s string) *string {
	return &s
}
