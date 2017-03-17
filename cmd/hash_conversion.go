package cmd

import "fmt"

func formatHashRate(hashrate uint64) string {
	prefix := []string{"", "K", "M", "G", "T", "P", "E", "Z", "Y"}
	value := float64(hashrate)
	idx := 0

	for value >= 1000 && idx < len(prefix) {
		idx++
		value /= 1000
	}

	return fmt.Sprintf("%.2f %sH/s", value, prefix[idx])
}
