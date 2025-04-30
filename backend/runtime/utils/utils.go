package utils

import (
	"fmt"
	"strconv"
)

// StrToFloat safely converts a str to float32
func StrToFloat(str string) (float32, error) {
	val, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0.0, fmt.Errorf("error converting string %s to float32: %q", str, err)
	}

	return float32(val), nil
}

// FloatToStr safely converts a float32 to str
func FloatToStr(val float32) string {
	return strconv.FormatFloat(float64(val), 'f', 2, 32)
}
