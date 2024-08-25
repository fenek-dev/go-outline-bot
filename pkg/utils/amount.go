package utils

import "strconv"

func ParseAmount(text string) (uint64, error) {
	amount, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
