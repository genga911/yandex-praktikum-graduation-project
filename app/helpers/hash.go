package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strconv"
)

func MakeMD5(in string) string {
	binHash := md5.Sum([]byte(in))
	return hex.EncodeToString(binHash[:])
}

func LuhnAlgorithm(digit string) bool {
	re := regexp.MustCompile(`[^\d]+`)
	digit = re.ReplaceAllString(digit, "")

	sum := 0
	digits := []rune(digit)
	j := len(digits)
	for i := 0; i < j; i++ {
		val, _ := strconv.Atoi(string(digits[i]))
		if i%2 != 0 {
			val = val * 2
			if val > 9 {
				val -= 9
			}
		}
		sum += val
	}

	return sum%10 == 0
}
