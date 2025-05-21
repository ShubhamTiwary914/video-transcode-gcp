package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const bucketName string = `temp-bucket-d08bdb88`

// region methods
// ====================

func main() {
	ctx := context.Background()
	cli, err := storage.NewClient(ctx)
	defer cli.Close()
	checkErr(err)
	bkt := cli.Bucket(bucketName)

	listGCS(ctx, bkt)
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
func uploadGCS(bkt *storage.BucketHandle, filepath string) {

}

// region utils
// ======================
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
