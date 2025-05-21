package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

var HLSbucket string

// region methods
// ====================

func main() {
	godotenv.Load()
	HLSbucket = os.Getenv("HLS_BUCKETNAME")

	ctx := context.Background()
	cli, err := storage.NewClient(ctx)
	defer cli.Close()
	checkErr(err)
	// bkt := cli.Bucket(HLSbucket)

	//worker-coroutines (background offloads)
	watcher, err := fsnotify.NewWatcher()
	checkErr(err)
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event.Name, " ", event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./out/stream_0")
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

func worker_routine(ctx context.Context, bkt *storage.BucketHandle, streamDir string) {
}

// list gcs objects in bucket
func listGCS(ctx context.Context, bkt *storage.BucketHandle) {
	it := bkt.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		checkErr(err)
		fmt.Println(attrs.Name, attrs.ContentType)
	}
}

// upload file from local to gcs bucket
func uploadGCS(ctx context.Context, bkt *storage.BucketHandle, file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}
	obj := bkt.Object(filepath.Base(file))
	w := obj.NewWriter(ctx)
	w.ChunkSize = 0 // no chunk upload (for small files)

	if _, err := w.Write(data); err != nil {
		log.Printf("upload error: %v", err)
	}
	if err := w.Close(); err != nil {
		log.Printf("writer close error: %v", err)
	}
}

// region utils
// ======================
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
