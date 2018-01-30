package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestLE64(t *testing.T) {
	var test []string
	test[0] = "test"
	out := hex.EncodeToString(PAE(test))
	fmt.Println(out)
}
