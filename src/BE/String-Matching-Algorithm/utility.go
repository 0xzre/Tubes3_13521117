package Algorithm

func LongestCommonSubstring(strA string, strB string) int {
	lengthA := len(strA)
	lengthB := len(strB)
	dp := make([][]int, lengthA+1)

	maxLen := 0
	for i := 0; i < lengthA+1; i++ {
		dp[i] = make([]int, lengthB+1)
		for j := 0; j < lengthB+1; j++ {
			if i == 0 || j == 0 {
				dp[i][j] = 0
			} else if strA[i-1] == strB[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				maxLen = max(maxLen, dp[i][j])
			} else {
				dp[i][j] = 0
			}
		}
	}

	percentage := float64(maxLen) / float64(lengthB) * 100
	return int(percentage)
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
