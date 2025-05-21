package main

import (
	// "bufio"
	"context"
	"fmt"
	"os/exec"

	// "fmt"
	"log"
	"os"
	"path/filepath"
	Utils "processjob/utils"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	// "github.com/joho/godotenv"
	// "google.golang.org/api/iterator"
)

var HLSbucket string
var ramfs_path string
var logs_path string

const DEBUG_MODE bool = false
const channelBufferSize int8 = 100
const streams int = 3

// region methods
// ====================

func main() {
	godotenv.Load()
	HLSbucket = os.Getenv("HLS_BUCKETNAME")
	ramfs_path = os.Getenv("RAMFS_PATH")
	logs_path = os.Getenv("LOGS_PATH")

	loggers := make([]*Utils.LogWriter, streams)
	Utils.SetupDirs(streams, ramfs_path, logs_path)
	Utils.InitLoggers(loggers, streams, channelBufferSize)

	ctx := context.Background()
	cli, err := storage.NewClient(ctx)
	defer cli.Close()
	checkErr(err)

	// worker-coroutines (background offloads)
	watchers := make([]*fsnotify.Watcher, streams)
	for i := 0; i < streams; i++ {
		watchers[i], err = fsnotify.NewWatcher()
		checkErr(err)
		defer watchers[i].Close()
		go GCS_offloader(watchers[i], loggers, i)
		err := watchers[i].Add(fmt.Sprintf("%sstream_%d", ramfs_path, i))
		checkErr(err)
	}

	cmd := exec.Command("bash", "./transcoder.sh", "./../../assets/trimmed_20.mp4")
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// offloads ffmpeg -> (ramfs/tmpfs) -> GCS bucket
func GCS_offloader(watcher *fsnotify.Watcher, loggers []*Utils.LogWriter, index int) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event.Name, " ", event.Op)
			loggers[index].Ch <- string(event.Name + " " + event.Op.String())

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// upload file from local to gcs bucket
func GCS_uploader(ctx context.Context, bkt *storage.BucketHandle, file string) {
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
