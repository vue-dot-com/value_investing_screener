package utils

import (
	"fmt"
	"strconv"
)

func CalculateTenCap(ownerearnings string) string {

	tenCapString := ""

	if ownerearnings != "" {
		n, _ := strconv.ParseFloat(ownerearnings, 32)
		tenCap := n * 10
		tenCapString := fmt.Sprintf("%f", tenCap)
		return tenCapString
	}

	return tenCapString
}
