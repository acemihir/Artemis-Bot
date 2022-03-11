package utils

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type GoogleFirebase struct {
	Context   context.Context
	App       *firebase.App
	Firestore *firestore.Client
}

var Firebase *GoogleFirebase

func SetupFirebase(file string) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(file)

	app, ex := firebase.NewApp(ctx, nil, sa)

	if ex != nil {
		log.Fatalf("[ERROR] Firebase app instantiation failed: %s", ex)
	}

	// Setup firestore
	fsclient, ex := app.Firestore(ctx)
	if ex != nil {
		log.Fatalf("[ERROR] Firestore app instantiation failed: %s", ex)
	}

	Firebase = &GoogleFirebase{
		Context:   ctx,
		App:       app,
		Firestore: fsclient,
	}
}

func (at *GoogleFirebase) SetFirestore(collection, doc string, data interface{}) {
	_, ex := at.Firestore.Collection(collection).Doc(doc).Set(at.Context, data)
	if ex != nil {
		log.Fatalf("[ERROR] Set in firestore failed: %s", ex)
	}
}

// Will return empty when no data is stored
func (at *GoogleFirebase) GetFirestore(collection, doc string) map[string]interface{} {
	dsnap, ex := at.Firestore.Collection(collection).Doc(doc).Get(at.Context)
	if ex != nil {
		if strings.Contains(ex.Error(), "not found") {
			return map[string]interface{}{}
		} else {
			log.Fatalf("[ERROR] Get from firestore failed: %s", ex)
		}
	}
	return dsnap.Data()
}

func (at *GoogleFirebase) DelFirestore(collection, doc string) {
	_, ex := at.Firestore.Collection(collection).Doc(doc).Delete(at.Context)
	if ex != nil {
		log.Fatalf("[ERROR] Delete from firestored failed: %s", ex)
	}
}
