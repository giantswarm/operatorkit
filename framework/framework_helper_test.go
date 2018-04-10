package framework

func contains(xs []string, x string) bool {
	for _, s := range xs {
		if s == x {
			return true
		}
	}

	return false
}
