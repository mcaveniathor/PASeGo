package past

import (
	"encoding/base64"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	test := []string{
		"Test Number One",
		"Test Number Two",
		"Three",
		"dsjflkjdlkjsd",
	}
	for _, str := range test {
		expected := base64.RawURLEncoding.EncodeToString([]byte(str))
		found, err := encodeConstantTime(str)
		if err != nil {
			t.Error(err)
		}
		if !(expected == found) {
			t.Errorf("Expected %s, got %s", expected, found)
		}
	}
}
