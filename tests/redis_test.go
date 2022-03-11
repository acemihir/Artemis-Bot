package utils

import (
	"testing"
	"time"

	"github.com/jerskisnow/Suggestions/src/utils"
)

func TestRedisSetup(t *testing.T) {
	utils.SetupCache(time.Duration(480) * time.Minute)
}

func TestRedisSet(t *testing.T) {
	utils.Cache.SetCache("test", "some value")
}

func TestRedisExists(t *testing.T) {
	res := utils.Cache.ExistsCache("test")
	if res == 0 {
		t.Fatal("Entry was not found.")
	}
}

func TestRedisGet(t *testing.T) {
	res := utils.Cache.GetCache("test")

	if res != "some value" {
		t.Fatal("Incorrect value returned.")
	}
}

func TestRedisDel(t *testing.T) {
	utils.Cache.DelCache("test")

	res := utils.Cache.ExistsCache("test")
	if res == 1 {
		t.Fatal("Entry was not deleted.")
	}
}
