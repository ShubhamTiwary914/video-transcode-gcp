package processjob

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	Utils "processjob/utils"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
)

// offloads ffmpeg -> (tmpfs/tmpfs) -> GCS bucket
func GCS_offloader(watcher *fsnotify.Watcher, loggers []*Utils.LogWriter, streamID int, fileID string, processedCtr map[int]int, uploadCh chan Utils.UploadEvent) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			loggers[streamID].Ch <- string(event.Name + " " + event.Op.String())
			//start upload after next CREATE
			if event.Op.String() == "CREATE" {
				dir, _, ext := Utils.GetFilePath_Split(event.Name)
				targetFile := fmt.Sprintf("%s%04d.ts", dir, processedCtr[streamID]-1)
				if ext == "m3u8" || processedCtr[streamID] > 0 {
					uploadCh <- Utils.UploadEvent{FilePath: targetFile, StreamID: streamID, FileID: fileID}
				}
				if ext != "m3u8" {
					processedCtr[streamID]++
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// upload files to HLS bucket
func GCS_uploader(ctx context.Context, bkt *storage.BucketHandle, localFile string, streamID int, fileID string) {
	start := time.Now()
	data, err := os.ReadFile(localFile)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}
	//filepath as it will appear in GCS bucket
	filebase := filepath.Base(localFile)
	gsFile := fmt.Sprintf("%s/%s/%s", fileID, Utils.StreamResolutions[streamID], filebase)
	//root of tmpfs (for master playlist)
	if streamID == -1 {
		gsFile = fmt.Sprintf("%s/%s", fileID, filebase)
	}
	//write to HLS-bucket
	obj := bkt.Object(gsFile)
	w := obj.NewWriter(ctx)
	w.ChunkSize = 0 // no chunk upload (better for small files)
	//delete local copy after upload
	if err = os.Remove(localFile); err != nil {
		log.Fatalf("failed deleting file: %v", err)
	}
	if _, err := w.Write(data); err != nil {
		log.Fatalf("upload error: %v", err)
		panic(err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("writer close error: %v", err)
		panic(err)
	}
	end := time.Since(start)
	log.Printf("Local: %s\tgsPath:%s\t (status: uploaded)  [%.2f secs]\n", localFile, gsFile, end.Seconds())
}
