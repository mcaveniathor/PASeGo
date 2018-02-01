package past

import (
	"bytes"
	"encoding/binary"
	"errors"
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

//decode6Bits decodes a 6 bit integer into an 8 bit integer
//using bitwise operators rather than a table
func decode6Bits(src int) byte {
	ret := -1
	ret += (((0x40 - src) & (src - 0x5b)) >> 8) & (src - 64)
	ret += (((0x60 - src) & (src - 0x7b)) >> 8) & (src - 70)
	ret += (((0x2f - src) & (src - 0x3a)) >> 8) & (src + 5)
	ret += (((0x2c - src) & (src - 0x2e)) >> 8) & 63
	ret += (((0x5e - src) & (src - 0x60)) >> 8) & 64
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

func decodeConstantTime(src string) (string, error) {
	dest := new(bytes.Buffer)
	srcBytes := []byte(src)
	srcLen := len(srcBytes)
	if srcLen == 0 {
		return "", errors.New("empty input")
	}
	var e byte
	e = 0
	var i int
	for i = 0; i+4 <= srcLen; i += 4 {
		chunk := srcBytes[i : i+5]
		b0 := decode6Bits(int(chunk[0]))
		b1 := decode6Bits(int(chunk[1]))
		b2 := decode6Bits(int(chunk[2]))
		b3 := decode6Bits(int(chunk[3]))
		err := dest.WriteByte(((b0 << 2) | (b1 >> 4)) & 0xff)
		if err != nil {
			return "", err
		}
		err = dest.WriteByte(((b1 << 4) | (b2 >> 2)) & 0xff)
		if err != nil {
			return "", err
		}
		err = dest.WriteByte(((b2 << 6) | (b3)) & 0xff)
		e |= (b0 | b1 | b2 | b3) >> 8
	}
	if i < srcLen {
		chunk := srcBytes[i:]
		b0 := decode6Bits(int(chunk[0]))
		if i+2 < srcLen {
			b1 := decode6Bits(int(chunk[1]))
			b2 := decode6Bits(int(chunk[2]))
			err := dest.WriteByte(((b0 << 2) | (b1 >> 4)) & 0xff)
			if err != nil {
				return "", err
			}
			err = dest.WriteByte(((b1 << 4) | (b2 >> 2)) & 0xff)
			if err != nil {
				return "", nil
			}
			e |= (b0 | b1 | b2) >> 8
		} else if i+1 < srcLen {
			b1 := decode6Bits(int(chunk[1]))
			err := dest.WriteByte(((b0 << 2) | (b1 >> 4)) & 0xff)
			if err != nil {
				return "", nil
			}
			e |= (b0 | b1) >> 8
		}
		if e != 0 {
			return "", errors.New("unexpected character found")
		}
	}
	return dest.String(), nil
}
