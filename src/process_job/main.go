package main

import (
	"context"
	"fmt"
	"os/exec"

	"io"
	"log"
	"os"
	GCS "processjob/gcs"
	Types "processjob/types"
	Utils "processjob/utils"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

const DEBUG_MODE bool = false
const logchannel_BufferSize int8 = 100
const uploadchannel_bufferSize int8 = 100
const streams int = 3

// region methods
// ====================
func main() {
	var err error
	start := time.Now()
	Env := Types.TasksEnv{}
	Proc := Types.Processor{}
	Channels := Types.ChannelsContainer{}

	NewEnvs(&Env)
	if err := Utils.SetupDirs(streams, &Env); err != nil {
		log.Fatalf("SetupDirs failed: %v", err)
	}
	NewProcessor(&Proc, &Env)
	defer Proc.Cli.Close()

	Channels.Loggers = make([]*Types.LogWriter, streams)
	Channels.UploadCh = make(chan Types.UploadEvent, uploadchannel_bufferSize)
	Proc.ProcessedCtr = make(map[int]int, streams)

	Utils.InitLoggers(Channels.Loggers, streams, Env.LOGS_PATH, logchannel_BufferSize)

	//>start worker co-routines + main(transcoder FFMPEG process)
	Proc.Watchers = make([]*fsnotify.Watcher, streams)
	startCoroutines(&Env, &Channels, &Proc)
	for i := 0; i < streams; i++ {
		defer Proc.Watchers[i].Close()
	}
	cmd := exec.Command("bash", "./transcoder.sh", Env.INPUT_PATH)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	finalChecks(&Env, &Channels, &Proc)
	end := time.Since(start)
	log.Printf("Time taken: %.2f sec", end.Seconds())
}

func NewEnvs(Env *Types.TasksEnv) {
	godotenv.Load()
	Env.HLS_BUCKET = os.Getenv("HLS_BUCKETNAME")
	Env.TMPFS_PATH = os.Getenv("TMPFS_PATH")
	Env.LOGS_PATH = os.Getenv("LOGS_PATH")
	Env.FILE_ID = os.Getenv("FILE_ID")
	Env.INPUT_PATH = os.Getenv("INPUT_PATH")
	Env.OUT_PATH = os.Getenv("OUT_PATH")
}

func NewProcessor(Proc *Types.Processor, Env *Types.TasksEnv) {
	var err error
	Proc.Ctx = context.Background()
	Proc.Cli, err = storage.NewClient(Proc.Ctx)
	Proc.Bkt = Proc.Cli.Bucket(Env.HLS_BUCKET)
	checkErr(err)
	//stdout to logfile
	logFile, err := os.OpenFile(fmt.Sprintf("%s/out.log", Env.OUT_PATH), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	checkErr(err)
	multi := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multi)
}

func finalChecks(Env *Types.TasksEnv, Channels *Types.ChannelsContainer, Proc *Types.Processor) {
	log.Println("\n\nDone with ffmpeg execution...")
	log.Println("Waiting for closing channel & remaining uploads (for mpegts - .ts files)...")
	close(Channels.UploadCh)
	Proc.UploadWg.Wait()
	log.Println("\nUploading remaining playlists...")
	uploadPlaylists(Env, Proc)
	fmt.Println("\nDone... ")
}

func startCoroutines(Env *Types.TasksEnv, Channels *Types.ChannelsContainer, Proc *Types.Processor) {
	var err error
	for i := 0; i < streams; i++ {
		Proc.Watchers[i], err = fsnotify.NewWatcher()
		checkErr(err)
		go logWorker(Channels.Loggers[i])
		go GCS.GCS_offloader(Env, Channels, Proc, i)
		Proc.UploadWg.Add(1)
		go uploadWorker(Channels.UploadCh, Proc, Env)
		err := Proc.Watchers[i].Add(fmt.Sprintf("%s/stream_%d", Env.TMPFS_PATH, i))
		checkErr(err)
	}
}

// Log receive & write routine
func logWorker(lw *Types.LogWriter) {
	for msg := range lw.Ch {
		lw.File.WriteString(msg + "\n")
	}
}

// Receive files & upload em to HLS-bucket
func uploadWorker(UploadCh <-chan Types.UploadEvent, Proc *Types.Processor, Env *Types.TasksEnv) {
	defer Proc.UploadWg.Done()
	for ev := range UploadCh {
		GCS.GCS_uploader(Proc, Env, ev.FilePath, ev.StreamID)
	}
}

// Clear backlog to upload remaining playlists
func uploadPlaylists(Env *Types.TasksEnv, Proc *Types.Processor) {
	//stream playlists
	for streamID := 0; streamID < streams; streamID++ {
		localfile := fmt.Sprintf("%s/stream_%d/playlist.m3u8", Env.TMPFS_PATH, streamID)
		GCS.GCS_uploader(Proc, Env, localfile, streamID)
	}
	//master
	GCS.GCS_uploader(Proc, Env, fmt.Sprintf("%s/master.m3u8", Env.TMPFS_PATH), -1)
}

// region utils
// ======================
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
