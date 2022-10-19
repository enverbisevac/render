package render

func totalPages(size, total int) int {
	quotient, remainder := total/size, total%size
	switch {
	case quotient == 0:
		return 1
	case remainder == 0:
		return quotient
	default:
		return quotient + 1
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if y == 0 {
		return x
	}
	if x < y {
		return x
	}
	return y
}
