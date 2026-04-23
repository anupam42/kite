package utils

import "strings"

// FuzzyScore returns a relevance score for pattern matched against str.
// Returns 0 if pattern chars do not appear in str as an ordered subsequence.
// Higher scores indicate better matches.
//
// Scoring tiers:
//   - Exact match:     1000
//   - Prefix match:    800–899
//   - Substring match: 500–699
//   - Fuzzy match:     1–499 (subsequence with bonuses for consecutive runs and word boundaries)
func FuzzyScore(pattern, str string) int {
	if pattern == "" || str == "" {
		return 0
	}
	p := strings.ToLower(pattern)
	s := strings.ToLower(str)

	if p == s {
		return 1000
	}
	if strings.HasPrefix(s, p) {
		// Shorter target scores higher within the prefix tier.
		return max(800, 899-len(s))
	}
	if idx := strings.Index(s, p); idx >= 0 {
		// Earlier occurrence and shorter target score higher within the substring tier.
		return max(500, 699-idx*3-len(s)/5)
	}
	return fuzzySubsequence(p, s)
}

// fuzzySubsequence checks whether all chars of pattern appear in str in order
// and returns a score based on match quality. Returns 0 if not a subsequence.
func fuzzySubsequence(pattern, str string) int {
	if len(pattern) > len(str) {
		return 0
	}
	pi := 0
	score := 0
	consecutive := 0
	lastMatch := -1

	for si := 0; si < len(str) && pi < len(pattern); si++ {
		if pattern[pi] != str[si] {
			consecutive = 0
			continue
		}
		charScore := 10
		// Word-boundary bonus: match at start of a segment.
		if si == 0 || str[si-1] == '-' || str[si-1] == '_' || str[si-1] == '.' {
			charScore += 30
		}
		// Consecutive-run bonus: adjacent matches score exponentially better.
		if lastMatch == si-1 {
			consecutive++
			charScore += consecutive * 15
		} else {
			consecutive = 0
		}
		score += charScore
		lastMatch = si
		pi++
	}
	if pi < len(pattern) {
		return 0
	}
	return score
}
