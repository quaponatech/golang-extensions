package convert

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// HexBytesToUpperString converts a byte array of hex numbers to a string
func HexBytesToUpperString(hexBytes []byte) string {
	if len(hexBytes) == 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(hexBytes))
}

// StringToHexBytes converts string to a byte array of hex numbers
func StringToHexBytes(hexString string) ([]byte, error) {
	if hexString == "" {
		return nil, fmt.Errorf("Empty string")
	}
	return hex.DecodeString(hexString)
}
