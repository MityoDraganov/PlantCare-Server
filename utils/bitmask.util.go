package utils

const (
	Monday    = 1 << iota // 1 (0000001)
	Tuesday               // 2 (0000010)
	Wednesday             // 4 (0000100)
	Thursday              // 8 (0001000)
	Friday                // 16 (0010000)
	Saturday              // 32 (0100000)
	Sunday                // 64 (1000000)
)


// ParseBitmask converts a bitmask (uint8) to a slice of integers (days of the week)
func ParseBitmask(bitmask uint8) []int {
	days := []int{}
	if bitmask&Monday != 0 {
		days = append(days, 1) // Monday
	}
	if bitmask&Tuesday != 0 {
		days = append(days, 2) // Tuesday
	}
	if bitmask&Wednesday != 0 {
		days = append(days, 3) // Wednesday
	}
	if bitmask&Thursday != 0 {
		days = append(days, 4) // Thursday
	}
	if bitmask&Friday != 0 {
		days = append(days, 5) // Friday
	}
	if bitmask&Saturday != 0 {
		days = append(days, 6) // Saturday
	}
	if bitmask&Sunday != 0 {
		days = append(days, 7) // Sunday
	}
	return days
}
