package past

type SymmetricKey struct {
	key string
}

//New Creates a SymmetricKey
func New(keyMaterial string) SymmetricKey {
	sk := SymmetricKey{key: keyMaterial}
	return sk
}

//encode encodes the SymmetricKey's
//key value using constant time base64url
func (s *SymmetricKey) encode() (string, error) {
	return encodeConstantTime(s.key)
}

func fromEncodedString(encoded string) SymmetricKey {
	return &SymmetricKey{key: decodeConstantTime(encoded)}
}
