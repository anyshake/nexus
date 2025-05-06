package main

func isLeap(y int) bool {
	return (y%400 == 0) || (y%4 == 0 && y%100 != 0)
}

func ldoy(y, m int) int {
	doy := [...]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365}
	if isLeap(y) && m >= 3 {
		return doy[m-1] + 1
	}
	return doy[m-1]
}

func mdy2dy(month, day, year int) int {
	return ldoy(year, month) + day - 1
}
