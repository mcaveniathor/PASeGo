package past

import "fmt"

//TypeError prints the expected and actual types in the event of a type mismatch
type TypeError struct {
	expected string
	found    string
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("Expected %q, found %q\n", t.expected, t.found)
}
