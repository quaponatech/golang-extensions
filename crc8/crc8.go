package crc8

import (
	sigurn_crc8 "github.com/sigurn/crc8"
)

// CRC8saeJ1850 calculates the crc for a given message by algorithm definition CRC-8/SAE-J1850
func CRC8saeJ1850(hexBytes []byte) byte {
	/*
		//crc ^= 0xff;
		crc := initialByte ^ 0xff
		//while (len--) {
		for _, value := range hexBytes {
			//crc ^= *data++;
			crc = crc ^ value
			//for (unsigned k = 0; k < 8; k++)
			for bit := 0; bit < 8; bit++ {
				//crc = crc & 0x80 ? (crc << 1) ^ 0x1d : crc << 1;
				if crc&byte(0x80) == byte(0x80) {
					crc = ((crc & byte(0x7f)) << 1) ^ byte(0x1d)
				} else {
					crc = (crc << 1)
				}
			}
		}
		//crc &= 0xff;
		crc = crc & 0xff
		//crc ^= 0xff;
		crc = crc ^ 0xff
		return crc
	*/
	algoDef := sigurn_crc8.MakeTable(sigurn_crc8.Params{
		Poly:   0x1D,
		Check:  0x00,
		Init:   0xFF,
		RefIn:  false,
		RefOut: false,
		XorOut: 0xFF,
		Name:   "CRC-8/SAE-J1850",
	})
	return sigurn_crc8.Checksum(hexBytes, algoDef)
}
