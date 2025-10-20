package utils

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/util/grand"
)

func GenRandomString(length int, validChars string) string {
	b := make([]string, length)
	chars := strings.Split(validChars, "")
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return strings.Join(b, "")
}

func SanitizeFilename(filename string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"..", "",
	)
	return replacer.Replace(strings.ReplaceAll(filename, "\u0000", ""))
}

func FindFirstMatch(s, reg string) string {
	re, err := regexp.Compile(reg)
	if err != nil {
		return ""
	}
	matches := re.FindStringSubmatch(s)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func ParseChineseNumberToInt(s string) int {
	if s == "" {
		return 0
	}

	clean := strings.TrimSpace(s)
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.ReplaceAll(clean, "，", "")

	if clean == "" {
		return 0
	}

	re := regexp.MustCompile(`^[\s]*([+-]?\d+(?:\.\d+)?)(?:\s*(万|亿))?`)
	m := re.FindStringSubmatch(clean)
	if len(m) == 0 {
		return 0
	}

	numStr := m[1]
	unit := m[2]

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	multiplier := 1.0
	switch unit {
	case "万":
		multiplier = 1e4
	case "亿":
		multiplier = 1e8
	}

	val = val * multiplier

	floored := math.Floor(val)

	if floored < float64(math.MinInt64) || floored > float64(math.MaxInt64) {
		return 0
	}

	return int(floored)
}

func RandomSecondsBatesInt(lower, upper int, n int) int {
	if n <= 0 {
		n = 1
	}
	if lower >= upper {
		return lower
	}

	sum := 0
	for i := 0; i < n; i++ {
		sum += grand.N(lower, upper)
	}
	iv := sum / n

	if iv < lower {
		iv = lower
	} else if iv > upper {
		iv = upper
	}
	return iv
}
