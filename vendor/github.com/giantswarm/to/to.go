package to

import "time"

func Duration(d *time.Duration) time.Duration {
	return *d
}

func DurationP(d time.Duration) *time.Duration {
	return &d
}

func Int(i *int) int {
	return *i
}

func IntP(i int) *int {
	return &i
}

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
