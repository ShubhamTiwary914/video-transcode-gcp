package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"log"
	"os"
	"path/filepath"
	Utils "processjob/utils"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

var HLSbucket string
var tmpfs_path string
var logs_path string
var fileID string
var inputPath string

const DEBUG_MODE bool = false
const logchannel_BufferSize int8 = 100
const uploadchannel_bufferSize int8 = 100
const streams int = 3

var ctx context.Context
var bkt *storage.BucketHandle
var uploadWg sync.WaitGroup

// region methods
// ====================
func main() {
	godotenv.Load()
	HLSbucket = os.Getenv("HLS_BUCKETNAME")
	tmpfs_path = os.Getenv("TMPFS_PATH")
	logs_path = os.Getenv("LOGS_PATH")
	fileID = os.Getenv("FILE_ID")
	inputPath = os.Getenv("INPUT_PATH")

	loggers := make([]*Utils.LogWriter, streams)
	uploadCh := make(chan Utils.UploadEvent, uploadchannel_bufferSize)
	processedCtr := make(map[int]int, streams)

	Utils.SetupDirs(streams, tmpfs_path, logs_path)
	Utils.InitLoggers(loggers, streams, logs_path, logchannel_BufferSize)

	ctx = context.Background()
	cli, err := storage.NewClient(ctx)
	defer cli.Close()
	checkErr(err)
	bkt = cli.Bucket(HLSbucket)

	//> worker-coroutines (background offloads)
	watchers := make([]*fsnotify.Watcher, streams)
	for i := 0; i < streams; i++ {
		watchers[i], err = fsnotify.NewWatcher()
		checkErr(err)
		defer watchers[i].Close()
		go logWorker(loggers[i])
		go GCS_offloader(watchers[i], loggers, i, processedCtr, uploadCh)
		uploadWg.Add(1)
		go uploadWorker(uploadCh)
		err := watchers[i].Add(fmt.Sprintf("%s/stream_%d", tmpfs_path, i))
		checkErr(err)
	}

	cmd := exec.Command("bash", "./transcoder.sh", inputPath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\n\nDone with ffmpeg execution...")
	log.Println("Waiting for closing channel & remaining uploads (for mpegts - .ts files)...")
	close(uploadCh)
	uploadWg.Wait()
	log.Println("Uploading remaining playlists...")
	uploadPlaylists(tmpfs_path)
	fmt.Println("Done!")
}

// Log receive & write routine
func logWorker(lw *Utils.LogWriter) {
	for msg := range lw.Ch {
		lw.File.WriteString(msg + "\n")
	}
}

// Receive files & upload em to HLS-bucket
func uploadWorker(uploadCh <-chan Utils.UploadEvent) {
	defer uploadWg.Done()
	for ev := range uploadCh {
		GCS_uploader(ev.FilePath, ev.StreamID, ev.FileID)
	}
}

// Clear backlog to upload remaining playlists
func uploadPlaylists(tmpfs_path string) {
	//stream playlists
	for streamID := 0; streamID < streams; streamID++ {
		GCS_uploader(fmt.Sprintf("%s/stream_%d/playlist.m3u8", tmpfs_path, streamID), streamID, fileID)
	}
	//master
	GCS_uploader(fmt.Sprintf("%s/master.m3u8", tmpfs_path), -1, fileID)
}

// offloads ffmpeg -> (tmpfs/tmpfs) -> GCS bucket
func GCS_offloader(watcher *fsnotify.Watcher, loggers []*Utils.LogWriter, streamID int, processedCtr map[int]int, uploadCh chan Utils.UploadEvent) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			loggers[streamID].Ch <- string(event.Name + " " + event.Op.String())
			//upload
			if event.Op.String() == "CREATE" {
				dir, _, ext := Utils.GetFilePath_Split(event.Name)
				targetFile := fmt.Sprintf("%s%04d.ts", dir, processedCtr[streamID]-1)
				if ext == "m3u8" || processedCtr[streamID] > 0 {
					uploadCh <- Utils.UploadEvent{FilePath: targetFile, StreamID: streamID, FileID: fileID}
					// fmt.Printf("Created: %s\t Target: %s\n", event.Name, targetFile)
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
func GCS_uploader(localFile string, streamID int, fileID string) {
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

	obj := bkt.Object(gsFile)
	w := obj.NewWriter(ctx)
	w.ChunkSize = 0 // no chunk upload (better for small files)

	if _, err := w.Write(data); err != nil {
		log.Fatalf("upload error: %v", err)
		panic(err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("writer close error: %v", err)
		panic(err)
	}
	log.Printf("Local: %s\tgsPath:%s\t (status: uploaded)\n", localFile, gsFile)
}

// region utils
// ======================
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
