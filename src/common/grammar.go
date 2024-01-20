package common

import (
	"strings"
)

// SnakeCase converts a PascalCase string into pascal_case
func SnakeCase(s string) string {
	if len(s) == 0 {
		return ""
	}
	points := make(map[int]struct{})
	for i := range s {
		if i > 0 && string(s[i]) == strings.ToUpper(string(s[i])) {
			points[i] = struct{}{}
		}
	}
	if len(points) == 0 {
		return strings.ToLower(s)
	}
	res := ""
	i := 0
	for {
		if i >= len(s) {
			break
		}
		if _, ok := points[i]; ok {
			res += "_"
		}
		res += strings.ToLower(string(s[i]))
		i++
	}
	return res
}

func KebabCase(s string) string {
	return strings.ReplaceAll(SnakeCase(s), "_", "-")
}

func IndefArticle(s string) string {
	lower := strings.ToLower(s)
	switch lower[0] {
	case 'a', 'e', 'i', 'o', 'u':
		return "an"
	}
	return "a"
}

func IsPlural(s string) bool {
	return strings.HasSuffix(s, "s")
}

func Pluralize(s string) string {
	if strings.ToLower(s) == "alias" {
		return s + "es"
	} else if IsPlural(s) {
		return s
	} else if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}
