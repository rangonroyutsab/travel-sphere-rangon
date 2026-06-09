package utils

import (
	"fmt"
	"strings"
)

func FormatPopulation(population int) string {
	if population >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", float64(population)/1_000_000_000)
	} else if population >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(population)/1_000_000)
	} else if population >= 1_000 {
		return fmt.Sprintf("%.1fK", float64(population)/1_000)
	}
	return fmt.Sprintf("%d", population)
}

func JoinStrings(values []string) string {
	if len(values) == 0 {
		return "N/A"
	}
	return strings.Join(values, ", ")
}

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacements := map[string]string{
		"'": "",
		".": "",
		",": "",
		"(": "",
		")": "",
		"&": "and",
	}
	for old, newValue := range replacements {
		value = strings.ReplaceAll(value, old, newValue)
	}
	return strings.Join(strings.Fields(value), "-")
}
