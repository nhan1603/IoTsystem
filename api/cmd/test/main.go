package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	data := readBinaryWatch(9)
	fmt.Println(data)
}

func readBinaryWatch(turnedOn int) []string {
	ans := []string{}
	minute := []int{8 * 60, 4 * 60, 2 * 60, 60, 32, 16, 8, 4, 2, 1}
	backTrack(&ans, minute, 0, turnedOn, 0)
	sort.Strings(ans)
	return ans
}

func backTrack(ans *[]string, minute []int, start int, turnedOn int, totalMinute int) {
	// Base case: if turnedOn LEDs are selected
	if turnedOn == 0 {
		fmt.Println(totalMinute)
		calculatedMinute := minuteToString(totalMinute)
		if calculatedMinute != "" {
			*ans = append(*ans, calculatedMinute)
		}
		return
	}

	// Iterate through the minute array
	for key := start; key < len(minute); key++ {
		// Include the current minute value in totalMinute
		newCalculated := totalMinute + minute[key]
		// Call backtrack with the next key (key + 1) to avoid reusing the same element
		backTrack(ans, minute, key+1, turnedOn-1, newCalculated)
	}
}

// minuteToString converts minutes to a "H:MM" format
func minuteToString(totalMinutes int) string {
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	if hours >= 12 {
		return ""
	}
	return strconv.Itoa(hours) + ":" + fmt.Sprintf("%02d", minutes)
}
