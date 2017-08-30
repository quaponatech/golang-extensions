package convert

import (
	"fmt"
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestHexBytesToUpperString(t *testing.T) {

	// Test list:
	// ----------
	// + EmptyInput
	// + 4BitNumberInput
	// + 32BitNumberInput
	// + StringConvertedInput
	// + LowerCaseInput
	// + UpperCaseInput

	t.Run("EmptyInput", func(t *testing.T) {
		emptyHexBytes := []byte{}
		test.AssertThat(t, HexBytesToUpperString(emptyHexBytes), "")
	})

	t.Run("8BitNumberInput", func(t *testing.T) {
		hexBytes := []byte{0x01, 0x02, 0x03}
		test.AssertThat(t, HexBytesToUpperString(hexBytes), "010203")
	})

	t.Run("4BitNumberInput", func(t *testing.T) {
		hexBytes := []byte{0x1, 0x2, 0x3}
		expectedString := "010203" // It is implictly converted to 1 byte (8bit) hex numbers
		test.AssertThat(t, HexBytesToUpperString(hexBytes), expectedString)
	})

	t.Run("32BitNumberInput", func(t *testing.T) {
		hexBytes := []byte{0x0001, 0x0002, 0x0003}
		expectedString := "010203" // It is implictly converted to 1 byte (8bit) hex numbers
		test.AssertThat(t, HexBytesToUpperString(hexBytes), expectedString)
	})

	t.Run("StringConvertedInput", func(t *testing.T) {
		hexBytes := []byte("123")
		expectedString := "313233" // The characters of the string are interpreted as bytes, so 0 is 30, 1 is 31, and so on
		test.AssertThat(t, HexBytesToUpperString(hexBytes), expectedString)
	})

	t.Run("LowerCaseInput", func(t *testing.T) {
		hexBytes := []byte{0x0a, 0x09, 0x0f}
		test.AssertThat(t, HexBytesToUpperString(hexBytes), "0A090F")
	})

	t.Run("UpperCaseInput", func(t *testing.T) {
		hexBytes := []byte{0x0A, 0x09, 0x0F}
		test.AssertThat(t, HexBytesToUpperString(hexBytes), "0A090F")
	})
}

func TestStringToHexBytes(t *testing.T) {

	// Test list:
	// ----------
	// + EmptyInput
	// + NoHexInput
	// + OddHexInput
	// + EvenHexInput
	// + LowerCaseInput
	// + UpperCaseInput

	t.Run("EmptyInput", func(t *testing.T) {
		_, err := StringToHexBytes("")
		test.AssertThat(t, err, "Empty string", "contains")
	})

	t.Run("NoHexInput", func(t *testing.T) {
		_, err := StringToHexBytes("HelloWorld")
		test.AssertThat(t, err, "encoding/hex: invalid byte: U+0048 'H'", "contains")
	})

	t.Run("OddHexInput", func(t *testing.T) {
		_, err := StringToHexBytes("012")
		test.AssertThat(t, err, "encoding/hex: odd length hex string", "contains")
	})

	t.Run("EvenHexInput", func(t *testing.T) {
		hexBytes, err := StringToHexBytes("0123")
		test.AssertThat(t, err, nil)
		test.AssertThat(t, len(hexBytes), 2)
		test.AssertThat(t, fmt.Sprintf("%v", hexBytes), "[1 35]")
		test.AssertThat(t, fmt.Sprintf("%X", hexBytes), "0123")
		test.AssertThat(t, int(hexBytes[0]), 0x01)
		test.AssertThat(t, int(hexBytes[1]), 0x23)
	})

	t.Run("LowerCaseInput", func(t *testing.T) {
		hexBytes, err := StringToHexBytes("0b1c2d3e")
		test.AssertThat(t, err, nil)
		test.AssertThat(t, len(hexBytes), 4)
		test.AssertThat(t, fmt.Sprintf("%x", hexBytes), "0b1c2d3e")
		test.AssertThat(t, fmt.Sprintf("%X", hexBytes), "0B1C2D3E")
	})

	t.Run("UpperCaseInput", func(t *testing.T) {
		hexBytes, err := StringToHexBytes("0B1C2D3E")
		test.AssertThat(t, err, nil)
		test.AssertThat(t, len(hexBytes), 4)
		test.AssertThat(t, fmt.Sprintf("%x", hexBytes), "0b1c2d3e")
		test.AssertThat(t, fmt.Sprintf("%X", hexBytes), "0B1C2D3E")
	})
}
