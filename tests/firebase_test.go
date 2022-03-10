package utils

import (
	"testing"

	"github.com/jerskisnow/Suggestions/src/utils"
)

func TestFirebaseSetup(t *testing.T) {
	utils.SetupFirebase("../firebase-credentials.json")
}

func TestFirestoreSet(t *testing.T) {
	utils.Firebase.SetFirestore("submissions", "abcdef", map[string]interface{}{
		"guild_id":   "1234",
		"channel_id": "2345",
		"message_id": "3456",
		"upvotes":    6,
		"downvotes":  2,
	})
}

func TestFirestoreGet(t *testing.T) {
	data := utils.Firebase.GetFirestore("submissions", "abcdef")
	if len(data) != 5 {
		t.Fatalf("Expected data map length of 5, got %d", len(data))
	}
	if data["guild_id"] != "1234" {
		t.Fatalf("Expected guild_id to be 1234, got %s", data["guild_id"])
	}
	t.Logf("Document data: %v", data)
}

func TestFirestoreNotExists(t *testing.T) {

}

func TestFirestoreDel(t *testing.T) {

}
