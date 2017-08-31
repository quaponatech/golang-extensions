package convert

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// HexBytesToUpperString converts a byte array of hex numbers to a string
func HexBytesToUpperString(hexBytes []byte) string {
	if len(hexBytes) == 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(hexBytes))
}

// CharBytesToUpperString converts a byte array of char representations to a string
func CharBytesToUpperString(charBytes []byte) (string, error) {
	if len(charBytes) == 0 {
		return "", fmt.Errorf("Empty string")
	}
	decoded, err := hex.DecodeString(fmt.Sprintf("%X", charBytes))
	if err != nil {
		return "", err
	}
	result := string(decoded)
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	if !isAlpha(result) {
		return "", fmt.Errorf("Not only alphabetical characters")
	}
	return strings.ToUpper(result), nil
}

// StringToHexBytes converts string to a byte array of hex numbers
func StringToHexBytes(hexString string) ([]byte, error) {
	if hexString == "" {
		return nil, fmt.Errorf("Empty string")
	}
	return hex.DecodeString(hexString)
}
