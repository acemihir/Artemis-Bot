package utils

import (
	"testing"

	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
)

func TestRandomString(t *testing.T) {
	res := utils.RandomString("abc!", 10)

	t.Logf("Custom ID: %s", res)

	if len(res) != 14 {
		t.Fatalf("Expected a length of 14, got %d", len(res))
	}
}
