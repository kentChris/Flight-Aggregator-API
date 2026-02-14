package util

import (
	"fmt"
	"strings"
)

func ParseDurationStringToInt(durationStr string) int {
	var hours, minutes int
	fmt.Sscanf(durationStr, "%dh %dm", &hours, &minutes)
	return (hours * 60) + minutes
}

func FormatIDR(amount float64) string {
	str := fmt.Sprintf("%.0f", amount)

	var result []string
	length := len(str)

	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}

	return "Rp " + strings.Join(result, ".")
}
