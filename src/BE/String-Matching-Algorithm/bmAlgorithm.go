package Algorithm

func bmSearch(text, pattern string) int {
	last := buildLast(pattern)
	n := len(text)
	m := len(pattern)
	i := m - 1
	if i > n-1 {
		return -1
	}
	j := m - 1
	for i <= n-1 {
		if pattern[j] == text[i] {
			if j == 0 {
				return i
			} else {
				i--
				j--
			}
		} else {
			lo, ok := last[text[i]]
			if !ok {
				lo = -1
			}
			i += m - min(j, 1+lo)
			j = m - 1
		}
	}
	return -1
}

func buildLast(pattern string) map[byte]int {
	last := make(map[byte]int)
	for i := 0; i < len(pattern); i++ {
		last[pattern[i]] = i
	}
	return last
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
