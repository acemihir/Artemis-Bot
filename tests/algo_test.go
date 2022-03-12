package utils

import (
	"testing"

	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func TestCreateId(t *testing.T) {
	res := utils.CreateId("abc!", 10)

	t.Logf("Custom ID: %s", res)

	if len(res) != 14 {
		t.Fatalf("Expected a length of 14, got %d", len(res))
	}
}
