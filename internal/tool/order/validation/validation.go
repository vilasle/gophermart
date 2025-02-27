package validation

func ValidNumber(number int) bool {
	result := checksum(number / 10)

	return (number%10+result)%10 == 0
}

func checksum(number int) int {
	var r int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		r += cur
		number = number / 10
	}
	return r % 10
}
