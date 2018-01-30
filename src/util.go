package past

import (
	"bytes"
	"crypto/subtle"
	"encoding/base64"
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

//encodeUnpadded encodes input using the unpadded alternate base64 encoding
//as defined in RFC 4648
func encodeUnpadded(src string) string {
	bytes := []byte(src)
	s := base64.URLEncoding.WithPadding(-1).EncodeToString(bytes)
	return s
}

func decodeUnpadded(src string) string {
	s := base64.URLEncoding.WithPadding(-1).DecodeString(src)
	return s
}

func hashEquals(s1, s2 string) bool {
	b1 := []byte(s1)
	b2 := []byte(s2)
	return bool(subtle.ConstantTimeCompare(b1, b2))
}

func validateAndRemoveFooter(params ...string) (string, error) {
	switch len(params) {
	case 1:
		payload := params[0]
		return payload, nil
	case 2:
		payload := params[0]
		footer := params[1]
	default:
		return "", errors.New("Incorrect number of parameters")
	}
	footer = encodeUnpadded(footer)
	payloadLen := len(payload)
	footerLen := len(footer)
	trailingLen := payloadLen - (footerLen + 1)
	trailing := payload[trailingLen:]

}
