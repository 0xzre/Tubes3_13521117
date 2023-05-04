package Algorithm

func KMPSearch(text string, pattern string) int {
	m := len(pattern)
	if m == 0 {
		return 0
	}

	lps := make([]int, m)
	j := 0

	computeLPS(pattern, lps)

	for i := 0; i < len(text); {
		if pattern[j] == text[i] {
			i++
			j++
		}

		if j == m {
			return i - j
		} else if i < len(text) && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return -1
}

func computeLPS(pattern string, lps []int) {
	m := len(pattern)
	len := 0
	lps[0] = 0
	i := 1

	for i < m {
		if pattern[i] == pattern[len] {
			len++
			lps[i] = len
			i++
		} else {
			if len != 0 {
				len = lps[len-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
}
