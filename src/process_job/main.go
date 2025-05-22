package main

import (
	"context"
	"fmt"
	"os/exec"
	"sync"

	"io"
	"log"
	"os"
	GCS "processjob/gcs"
	Utils "processjob/utils"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

var (
	HLSbucket  string
	tmpfs_path string
	logs_path  string
	fileID     string
	inputPath  string
	outPath    string
)

var (
	ctx      context.Context
	bkt      *storage.BucketHandle
	uploadWg sync.WaitGroup
)

const DEBUG_MODE bool = false
const logchannel_BufferSize int8 = 100
const uploadchannel_bufferSize int8 = 100
const streams int = 3

// region methods
// ====================
func main() {
	start := time.Now()
	initialize()
	loggers := make([]*Utils.LogWriter, streams)
	uploadCh := make(chan Utils.UploadEvent, uploadchannel_bufferSize)
	processedCtr := make(map[int]int, streams)
	Utils.InitLoggers(loggers, streams, logs_path, logchannel_BufferSize)

	cli, err := storage.NewClient(ctx)
	defer cli.Close()
	checkErr(err)
	bkt = cli.Bucket(HLSbucket)

	//>start worker co-routines + main(transcoder FFMPEG process)
	watchers := make([]*fsnotify.Watcher, streams)
	startCoroutines(watchers, loggers, processedCtr, uploadCh)
	for i := 0; i < streams; i++ {
		defer watchers[i].Close()
	}
	cmd := exec.Command("bash", "./transcoder.sh", inputPath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	finalChecks(uploadCh)
	end := time.Since(start)
	log.Printf("Time taken: %.2f sec", end.Seconds())
}

func initialize() {
	godotenv.Load()
	HLSbucket = os.Getenv("HLS_BUCKETNAME")
	tmpfs_path = os.Getenv("TMPFS_PATH")
	logs_path = os.Getenv("LOGS_PATH")
	fileID = os.Getenv("FILE_ID")
	inputPath = os.Getenv("INPUT_PATH")
	outPath = os.Getenv("OUT_PATH")
	Utils.SetupDirs(streams, tmpfs_path, logs_path)
	ctx = context.Background()
	//stdout to logfile
	logFile, err := os.OpenFile(fmt.Sprintf("%s/out.log", outPath), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	checkErr(err)
	multi := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multi)

}

func finalChecks(uploadCh chan Utils.UploadEvent) {
	log.Println("\n\nDone with ffmpeg execution...")
	log.Println("Waiting for closing channel & remaining uploads (for mpegts - .ts files)...")
	close(uploadCh)
	uploadWg.Wait()
	log.Println("\nUploading remaining playlists...")
	uploadPlaylists(tmpfs_path)
	fmt.Println("\nDone... ")
}

func startCoroutines(watchers []*fsnotify.Watcher, loggers []*Utils.LogWriter, processedCtr map[int]int, uploadCh chan Utils.UploadEvent) {
	var err error
	for i := 0; i < streams; i++ {
		watchers[i], err = fsnotify.NewWatcher()
		checkErr(err)
		go logWorker(loggers[i])
		go GCS.GCS_offloader(watchers[i], loggers, i, fileID, processedCtr, uploadCh)
		uploadWg.Add(1)
		go uploadWorker(uploadCh)
		err := watchers[i].Add(fmt.Sprintf("%s/stream_%d", tmpfs_path, i))
		checkErr(err)
	}
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
		GCS.GCS_uploader(ctx, bkt, ev.FilePath, ev.StreamID, ev.FileID)
	}
}

// Clear backlog to upload remaining playlists
func uploadPlaylists(tmpfs_path string) {
	//stream playlists
	for streamID := 0; streamID < streams; streamID++ {
		GCS.GCS_uploader(ctx, bkt, fmt.Sprintf("%s/stream_%d/playlist.m3u8", tmpfs_path, streamID), streamID, fileID)
	}
	//master
	GCS.GCS_uploader(ctx, bkt, fmt.Sprintf("%s/master.m3u8", tmpfs_path), -1, fileID)
}

// region utils
// ======================
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
