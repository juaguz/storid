package currencies

import (
	"fmt"
	"strconv"
)

type Number interface {
	~float64 | ~int32 | ~float32
}

func FloatToCents[T Number](f T) int {
	return int(float64(int(f*100)) / 100.0 * 100)
}

func CentsToNumber[T Number](c int) T {
	return T(float64(c) / 100.0)
}

func StringToCents(s string) (int, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	f = float64(int(f*100)) / 100.0
	return FloatToCents(f), nil
}

func CentsToString(c int) string {
	return fmt.Sprintf("%.2f", float64(c)/100.0)
}
