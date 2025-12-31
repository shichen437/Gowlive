package utils

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

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

func RemoveEmoji(s string) string {
	out := make([]rune, 0, utf8.RuneCountInString(s))
	for _, r := range s {
		if isEmojiRune(r) {
			continue
		}
		if !unicode.IsPrint(r) {
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

func ReplaceColonWithDash(name string) string {
	return strings.ReplaceAll(name, ":", "-")
}

func isEmojiRune(r rune) bool {
	// 常见表情符号与符号的范围（不完全，但覆盖大部分）
	// Emoticons: U+1F600–U+1F64F
	if r >= 0x1F600 && r <= 0x1F64F {
		return true
	}
	// Misc Symbols and Pictographs: U+1F300–U+1F5FF
	if r >= 0x1F300 && r <= 0x1F5FF {
		return true
	}
	// Transport and Map Symbols: U+1F680–U+1F6FF
	if r >= 0x1F680 && r <= 0x1F6FF {
		return true
	}
	// Supplemental Symbols and Pictographs: U+1F900–U+1F9FF
	if r >= 0x1F900 && r <= 0x1F9FF {
		return true
	}
	// Symbols and Pictographs Extended-A: U+1FA70–U+1FAFF
	if r >= 0x1FA70 && r <= 0x1FAFF {
		return true
	}
	// Dingbats: U+2700–U+27BF
	if r >= 0x2700 && r <= 0x27BF {
		return true
	}
	// Miscellaneous Symbols: U+2600–U+26FF
	if r >= 0x2600 && r <= 0x26FF {
		return true
	}
	// Enclosed Alphanumeric Supplement: U+1F100–U+1F1FF
	if r >= 0x1F100 && r <= 0x1F1FF {
		return true
	}
	// Enclosed Ideographic Supplement: U+1F200–U+1F2FF
	if r >= 0x1F200 && r <= 0x1F2FF {
		return true
	}
	// Flags (Regional Indicator Symbols): U+1F1E6–U+1F1FF
	if r >= 0x1F1E6 && r <= 0x1F1FF {
		return true
	}
	// Variation Selectors: U+FE0F 等（常用于把字符变成 emoji 表现）
	if r == 0xFE0F || r == 0xFE0E {
		return true
	}
	// Zero Width Joiner（用于组合 emoji，去掉它能打散序列）
	if r == 0x200D {
		return true
	}
	return false
}
