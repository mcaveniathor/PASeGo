package past

import (
	"bytes"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"reflect"
)

//LE64 encodes a 64-bit unsigned integer into a little-endian binary string.
func LE64(i int) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint64(i))
	return buf.String()
}

//PAE returns the number of pieces (8 bytes),
//then each piece prefixed by its length(LE64 encoded)
func PAE(pieces []string) (string, error) {
	t := reflect.TypeOf(pieces).Name()
	if t != "[string]" {
		return "", &TypeError{
			expected: "[]string",
			found:    t,
		}
	}
	output := LE64(len(pieces))
	for _, piece := range pieces {
		output = output + LE64(len(piece)) + piece
	}
	return output, nil
}

func hashEquals(s1, s2 string) bool {
	b1 := []byte(s1)
	b2 := []byte(s2)
	return subtle.ConstantTimeCompare(b1, b2) == 1
}

func validateAndRemoveFooter(params ...string) (string, error) {
	var payload, encodedFooter string
	switch len(params) {
	case 1:
		payload = params[0]
		return payload, nil
	case 2:
		payload = params[0]
		footer := params[1]
		var err error
		encodedFooter, err = encodeConstantTime(footer)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("Incorrect number of parameters")
	}
	payloadLen := len(payload)
	footerLen := len(encodedFooter)
	trailingLen := payloadLen - (footerLen + 1)
	trailing := payload[trailingLen:]
	if !hashEquals(encodedFooter, trailing) {
		return "", errors.New("Invalid message footer")
	}
	return payload[:trailingLen], nil
}
