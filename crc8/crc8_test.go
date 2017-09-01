package crc8

import (
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestCRC8saeJ1850(t *testing.T) {
	actual := CRC8saeJ1850([]byte{0x00, 0x10, 0x10, 0x10, 0x00, 0x01})
	test.AssertThat(t, int(actual), 0x38)

	actual = CRC8saeJ1850([]byte{0x00, 0x10, 0x10, 0x1F, 0xFF, 0xFF})
	test.AssertThat(t, int(actual), 0x26)
}
