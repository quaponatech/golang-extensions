package crc8

// CRC8saeJ1850 calculates the crc for a given message over
func CRC8saeJ1850(hexBytes []byte, initialByte byte) byte {
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
}
