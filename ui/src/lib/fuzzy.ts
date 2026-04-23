/**
 * Score how well `pattern` matches `str`.
 * Returns 0 if the pattern characters do not appear in str as an ordered subsequence.
 * Higher scores mean better matches.
 *
 * Tiers:
 *   Exact match     → 1000
 *   Prefix match    → 800–899
 *   Substring match → 500–699
 *   Fuzzy match     → 1–499  (subsequence with consecutive-run and word-boundary bonuses)
 */
export function fuzzyScore(pattern: string, str: string): number {
  if (!pattern || !str) return 0

  const p = pattern.toLowerCase()
  const s = str.toLowerCase()

  if (p === s) return 1000
  if (s.startsWith(p)) return Math.max(800, 899 - s.length)

  const idx = s.indexOf(p)
  if (idx >= 0) return Math.max(500, 699 - idx * 3 - Math.floor(s.length / 5))

  return fuzzySubsequence(p, s)
}

function fuzzySubsequence(pattern: string, str: string): number {
  if (pattern.length > str.length) return 0

  let pi = 0
  let score = 0
  let consecutive = 0
  let lastMatch = -1

  for (let si = 0; si < str.length && pi < pattern.length; si++) {
    if (pattern[pi] !== str[si]) {
      consecutive = 0
      continue
    }

    let charScore = 10
    if (si === 0 || str[si - 1] === '-' || str[si - 1] === '_' || str[si - 1] === '.') {
      charScore += 30
    }
    if (lastMatch === si - 1) {
      consecutive++
      charScore += consecutive * 15
    } else {
      consecutive = 0
    }

    score += charScore
    lastMatch = si
    pi++
  }

  return pi < pattern.length ? 0 : score
}
