package main

import (
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// DatabaseCleanupHandler
func DatabaseCleanupHandler(w http.ResponseWriter, r *http.Request) {

	// Disable for now
	return

	//---------------------
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Context loaded. Starting execution.")

	firebasePID := os.Getenv("FIREBASE_PROJECT")
	log.Debugf(ctx, "Loaded firebase project ID: %v", firebasePID)
	if firebasePID == "" {
		log.Criticalf(ctx, "Firebase Project ID is nil, I cannot continue.")
		panic("Firebase Project ID is nil")
	}

	fbClient, err := firestore.NewClient(ctx, firebasePID)
	defer fbClient.Close()
	if err != nil {
		log.Errorf(ctx, "Could not create new client for Firebase %v", err)
		w.WriteHeader(500)
		return
	}

	dcIter := fbClient.Collection("departments").Documents(ctx)
	batch := fbClient.Batch()
	batchSize := 0
	for {
		doc, err := dcIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Error 35: %v", err)
			return
		}
		docsToDelete := doc.Ref.Collection("2018").Documents(ctx)
		for {
			inner, err := docsToDelete.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error 45: %v", err)
				return
			}
			if batchSize > 400 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return
				}
				batchSize = 0
				batch = fbClient.Batch()
			}
			batchSize++
			batch.Delete(inner.Ref)
		}
		docsToDelete = doc.Ref.Collection("courses").Documents(ctx)
		for {
			inner, err := docsToDelete.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error 66: %v", err)
				return
			}
			if batchSize > 400 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return
				}
				batchSize = 0
				batch = fbClient.Batch()
			}
			batchSize++
			batch.Delete(inner.Ref)
		}
		docsToDelete = doc.Ref.Collection("courses_fall").Documents(ctx)
		for {
			inner, err := docsToDelete.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error 87: %v", err)
				return
			}
			if batchSize > 400 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return
				}
				batchSize = 0
				batch = fbClient.Batch()
			}
			batchSize++
			batch.Delete(inner.Ref)
		}
		docsToDelete = doc.Ref.Collection("courses_spring").Documents(ctx)
		for {
			inner, err := docsToDelete.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error 108: %v", err)
				return
			}
			if batchSize > 400 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return
				}
				batchSize = 0
				batch = fbClient.Batch()
			}
			batchSize++
			batch.Delete(inner.Ref)
		}
		docsToDelete = doc.Ref.Collection("courses_summer").Documents(ctx)
		for {
			inner, err := docsToDelete.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Errorf(ctx, "Error 129: %v", err)
				return
			}
			if batchSize > 400 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return
				}
				batchSize = 0
				batch = fbClient.Batch()
			}
			batchSize++
			batch.Delete(inner.Ref)
		}

	}

}
