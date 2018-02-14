package past

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	//HEADER defines the paseto protocol version in use
	HEADER string = v2
)

//aeadEncrypt Encrypts a message given a shared key
func aeadEncrypt(data, header, footer, nonceTest string, sk *SymmetricKey) (string, error) {
	aead, err := chacha20poly1305.New([]byte(sk.key))
	if err != nil {
		return "", err
	}
	plaintext := []byte(data)
	var nonce []byte
	if nonceTest {
		nonce = []byte(nonceTest)
	} else {
		nonce = make([]byte, 24)
		_, err := io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return "", err
		}
		nonce, err = blake2b.New256(nonce)
		if err != nil {
			return "", err
		}
	}
	params := []string{header, nonce, footer}
	additionalData, err := PAE(params)
	cipherText := aead.Seal(nil, nonce, plaintext, additionalData)
	nc, err := encodeConstantTime(nonce + cipherText)
	if footer {
		if err != nil {
			return "", err
		}
		f, err := encodeConstantTime(footer)
		return header + nc + "." + f, nil
	}
	return header + nc, nil
}
