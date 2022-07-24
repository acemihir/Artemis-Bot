package utils

import (
	"testing"
	"time"

	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
)

func TestRedisSetup(t *testing.T) {
	utils.SetupCache(time.Duration(480) * time.Minute)
}

func TestRedisSet(t *testing.T) {
	ex := utils.Cache.SetCache("test", "some value")
	if ex != nil {
		t.Fatalf("Set in Redis failed: %v", ex)
	}
}

func TestRedisExists(t *testing.T) {
	res, ex := utils.Cache.ExistsCache("test")
	if ex != nil {
		t.Fatalf("Exists in Redis failed: %v", ex)
	}
	if res == 0 {
		t.Fatal("Entry was not found")
	}
}

func TestRedisGet(t *testing.T) {
	res, ex := utils.Cache.GetCache("test")
	if ex != nil {
		t.Fatalf("Get from Redis failed: %v", ex)
	}

	if res != "some value" {
		t.Fatal("Incorrect value returned")
	}
}

func TestRedisDel(t *testing.T) {
	utils.Cache.DelCache("test")

	res, ex := utils.Cache.ExistsCache("test")
	if ex != nil {
		t.Fatalf("Delete from Redis failed: %v", ex)
	}
	if res == 1 {
		t.Fatal("Entry was not deleted")
	}
}
