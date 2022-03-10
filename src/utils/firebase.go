package utils

import (
	"context"
	"log"

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
		log.Fatalf("Could not instantiate firebase app: %s", ex)
	}

	// Setup firestore
	fsclient, ex := app.Firestore(ctx)
	if ex != nil {
		log.Fatalf("Could not instantiate firestore: %s", ex)
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
		log.Fatalf("Could not fetch from firestore: %s", ex)
	}
}

// Will return empty when no data is stored
func (at *GoogleFirebase) GetFirestore(collection, doc string) map[string]interface{} {
	dsnap, ex := at.Firestore.Collection(collection).Doc(doc).Get(at.Context)
	if ex != nil {
		log.Fatalf("Could not get from firestore: %s", ex)
	}
	return dsnap.Data()
}

func (at *GoogleFirebase) DelFirestore(collection, doc string) {
	_, ex := at.Firestore.Collection(collection).Doc(doc).Delete(at.Context)
	if ex != nil {
		log.Fatalf("Could not delete from firestore: %s", ex)
	}
}
