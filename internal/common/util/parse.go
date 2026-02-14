package util

import "fmt"

func ParseDurationStringToInt(durationStr string) int {
	var hours, minutes int
	fmt.Sscanf(durationStr, "%dh %dm", &hours, &minutes)
	return (hours * 60) + minutes
}
