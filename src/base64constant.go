package past

import (
	"bytes"
	"encoding/binary"
)

//encode6Bits encodes an 8 bit integer to 6 bits
//using bitwise operators rather than a table
func encode6Bits(src int) byte {
	diff := (0x41)
	diff += ((25 - src) >> 8) & 6
	diff -= ((51 - src) >> 8) & 75
	diff -= ((61 - src) >> 8) & 15
	diff += ((62 - src) >> 8) & 3
	ret := src + diff
	return byte(ret)
}

//encodeConstantTime base64 encodes a string
//without padding in constant time (timing-attack safe)
func encodeConstantTime(src string) (string, error) {
	dest := new(bytes.Buffer)
	srcBytes := []byte(src)
	srcLen := len(srcBytes) //since slice length is stored in golang, this is timing-attack safe
	//encode 3 bytes
	var i int
	for i = 0; i+3 <= srcLen; i += 3 {
		chunk := srcBytes[i : i+4]
		b0 := chunk[0]
		b1 := chunk[1]
		b2 := chunk[2]
		err := binary.Write(dest, binary.LittleEndian, encode6Bits(int(b0>>2)))
		if err != nil {
			return "", err
		}
		err = binary.Write(dest, binary.LittleEndian, encode6Bits(int(((b0<<4)|(b1>>4))&63)))
		if err != nil {
			return "", err
		}
		err = binary.Write(dest, binary.LittleEndian, encode6Bits(int(((b1<<2)|(b2>>6))&63)))
		if err != nil {
			return "", err
		}
		err = binary.Write(dest, binary.LittleEndian, encode6Bits(int(b2&63)))
		if err != nil {
			return "", err
		}
	}
	if i < srcLen {
		chunk := srcBytes[i:]
		b0 := chunk[0]
		if i+1 < srcLen {
			b1 := chunk[1]
			err := binary.Write(dest, binary.LittleEndian, encode6Bits(int(b0>>2)))
			if err != nil {
				return "", err
			}
			err = binary.Write(dest, binary.LittleEndian, encode6Bits(int(((b0<<4)|(b1>>4))&63)))
			if err != nil {
				return "", err
			}
			err = binary.Write(dest, binary.LittleEndian, encode6Bits(int((b1<<2))&63))
			if err != nil {
				return "", err
			}
		} else {
			err := binary.Write(dest, binary.LittleEndian, encode6Bits(int(b0>>2)))
			if err != nil {
				return "", err
			}
			err = binary.Write(dest, binary.LittleEndian, encode6Bits(int((b0<<4)&63)))
			if err != nil {
				return "", err
			}
		}
	}
	return dest.String(), nil

}
