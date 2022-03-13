package utils

import (
	"context"
	"os"
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
		Cout("[ERROR] Firebase app instantiation failed: %s", Red, ex)
		os.Exit(1)
	}

	// Setup Firestore
	fsclient, ex := app.Firestore(ctx)
	if ex != nil {
		Cout("[ERROR] Firestore app instantiation failed: %s", Red, ex)
		os.Exit(1)
	}

	Firebase = &GoogleFirebase{
		Context:   ctx,
		App:       app,
		Firestore: fsclient,
	}
}

func (at *GoogleFirebase) SetFirestore(collection, doc string, data interface{}) error {
	_, ex := at.Firestore.Collection(collection).Doc(doc).Set(at.Context, data)
	return ex
}

// Will return empty when no data is stored
func (at *GoogleFirebase) GetFirestore(collection, doc string) (map[string]interface{}, error) {
	dsnap, ex := at.Firestore.Collection(collection).Doc(doc).Get(at.Context)
	if ex != nil {
		if strings.Contains(ex.Error(), "not found") {
			return map[string]interface{}{}, nil
		} else {
			return map[string]interface{}{}, ex
		}
	}
	return dsnap.Data(), nil
}

func (at *GoogleFirebase) DelFirestore(collection, doc string) error {
	_, ex := at.Firestore.Collection(collection).Doc(doc).Delete(at.Context)
	return ex
}
