package past

//encodeConstantTime base64 encodes a string
//with or without padding in constant time
func encodeConstantTime(src string, bool pad) string {
	dest := ""
	srcBytes := []byte(src)
	srcLen := len(srcBytes)
	for i; i+3 <= srcLen; i += 3 {
		chunk := srcBytes[i : i+3]
		b0 := chunk[0]
		b1 := chunk[1]
		b2 := chunk[2]
		dest
	}

}
