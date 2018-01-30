package past

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
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
			found:    string(t),
		}
	}
	output := LE64(len(pieces))
	for _, piece := range pieces {
		output = output + LE64(len(piece)) + piece
	}
	return output, nil
}

//encodeUnpadded Encodes input using the unpadded alternate base64 encoding
//as defined in RFC 4648
func encodeUnpadded(src int) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, src)
	s := base64.URLEncoding.WithPadding(-1).EncodeToString(buf.Bytes())
	return s
}
