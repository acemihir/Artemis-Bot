package utils

import (
	"testing"

	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func TestFirebaseSetup(t *testing.T) {
	utils.SetupFirebase("../firebase-credentials.json")
}

func TestFirestoreSet(t *testing.T) {
	ex := utils.Firebase.SetFirestore("submissions", "abcdef", map[string]interface{}{
		"guild_id":   "1234",
		"channel_id": "2345",
		"message_id": "3456",
		"upvotes":    6,
		"downvotes":  2,
	})
	if ex != nil {
		t.Fatalf("Set in Firestore failed: %v", ex)
	}
}

func TestFirestoreGet(t *testing.T) {
	data, ex := utils.Firebase.GetFirestore("submissions", "abcdef")
	if ex != nil {
		t.Fatalf("Get from Firestore failed: %v", ex)
	}
	if len(data) != 5 {
		t.Fatalf("Expected data map length of 5, got %d", len(data))
	}
	if data["guild_id"] != "1234" {
		t.Fatalf("Expected guild_id to be 1234, got %s", data["guild_id"])
	}
	t.Logf("Document data: %v", data)
}

func TestFirestoreDel(t *testing.T) {
	ex := utils.Firebase.DelFirestore("submissions", "abcdef")
	if ex != nil {
		t.Fatalf("Delete form Firestore failed: %v", ex)
	}
}

func TestFirestoreNotExists(t *testing.T) {
	data, ex := utils.Firebase.GetFirestore("submissions", "abcdef")
	if ex != nil {
		t.Fatalf("Firestore get failed: %v", ex)
	}
	if len(data) != 0 {
		t.Fatalf("Expected data map length of 0, got %d", len(data))
	}
	t.Logf("Document data: %v", data)
}
