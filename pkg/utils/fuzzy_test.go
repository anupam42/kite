package utils

import "testing"

func TestFuzzyScore(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		str     string
		wantMin int // score must be >= wantMin
		wantMax int // score must be <= wantMax (0 means exact check)
	}{
		{name: "empty pattern", pattern: "", str: "nginx", wantMin: 0, wantMax: 0},
		{name: "empty string", pattern: "nginx", str: "", wantMin: 0, wantMax: 0},
		{name: "exact match", pattern: "nginx", str: "nginx", wantMin: 1000, wantMax: 1000},
		{name: "exact match case-insensitive", pattern: "NGINX", str: "nginx", wantMin: 1000, wantMax: 1000},
		{name: "prefix match", pattern: "nginx", str: "nginx-deployment", wantMin: 800, wantMax: 899},
		{name: "substring match", pattern: "api", str: "my-api-server", wantMin: 500, wantMax: 699},
		{name: "fuzzy subsequence", pattern: "nd", str: "nginx-deployment", wantMin: 1, wantMax: 499},
		{name: "no match", pattern: "xyz", str: "nginx", wantMin: 0, wantMax: 0},
		{name: "pattern longer than string", pattern: "toolong", str: "too", wantMin: 0, wantMax: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FuzzyScore(tt.pattern, tt.str)
			if tt.wantMax == 0 && tt.wantMin == 0 {
				if got != 0 {
					t.Fatalf("FuzzyScore(%q, %q) = %d, want 0", tt.pattern, tt.str, got)
				}
				return
			}
			if got < tt.wantMin || got > tt.wantMax {
				t.Fatalf("FuzzyScore(%q, %q) = %d, want [%d, %d]", tt.pattern, tt.str, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestFuzzyScoreOrdering(t *testing.T) {
	// Verify the relative ordering of match quality: exact > prefix > substring > fuzzy
	pattern := "api"
	exact := FuzzyScore(pattern, "api")
	prefix := FuzzyScore(pattern, "api-server")
	substring := FuzzyScore(pattern, "my-api-server")
	fuzzy := FuzzyScore(pattern, "a-pod-image")

	if exact <= prefix {
		t.Errorf("exact(%d) should beat prefix(%d)", exact, prefix)
	}
	if prefix <= substring {
		t.Errorf("prefix(%d) should beat substring(%d)", prefix, substring)
	}
	if substring <= fuzzy {
		t.Errorf("substring(%d) should beat fuzzy(%d)", substring, fuzzy)
	}
	if fuzzy == 0 {
		t.Errorf("fuzzy(%d) should be > 0 for subsequence match", fuzzy)
	}
}

func TestFuzzyScoreWordBoundaryBonus(t *testing.T) {
	// "nd" is not a substring of either target, so both fall through to fuzzySubsequence.
	// "nginx-deployment" has both chars at word boundaries (start and after '-');
	// "kubernetes-node" matches 'n' mid-word and 'd' mid-word — lower score expected.
	bothBoundary := FuzzyScore("nd", "nginx-deployment")
	noBoundary := FuzzyScore("nd", "kubernetes-node")
	if bothBoundary <= noBoundary {
		t.Errorf("boundary match (%d) should beat mid-word match (%d)", bothBoundary, noBoundary)
	}
}
