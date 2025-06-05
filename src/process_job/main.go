package main

import (
	"context"
	"fmt"
	"os/exec"

	"bytes"
	"io"
	"log"
	"os"
	GCS "processjob/gcs"
	Types "processjob/types"
	Utils "processjob/utils"
	"time"

	"cloud.google.com/go/compute/metadata"
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
	start := time.Now()
	Env := Types.TasksEnv{}
	Proc := Types.Processor{}
	Channels := Types.ChannelsContainer{}

	NewEnvs(&Env)
	if err := Utils.SetupDirs(streams, &Env); err != nil {
		log.Panicf("SetupDirs failed: %v", err)
	}
	fmt.Println("Directories Setup done!")
	NewProcessor(&Proc, &Env)
	defer Proc.Cli.Close()

	Channels.Loggers = make([]*Types.LogWriter, streams)
	Channels.UploadCh = make(chan Types.UploadEvent, uploadchannel_bufferSize)
	Proc.ProcessedCtr = make(map[int]int, streams)

	Utils.InitLoggers(Channels.Loggers, streams, Env.LOGS_PATH, logchannel_BufferSize)
	fmt.Println("Channels Initialized!")

	//>start worker co-routines + main(transcoder FFMPEG process)
	Proc.Watchers = make([]*fsnotify.Watcher, streams)
	startCoroutines(&Env, &Channels, &Proc)
	fmt.Println("Coroutines Started!")
	for i := 0; i < streams; i++ {
		defer Proc.Watchers[i].Close()
	}
	fmt.Println("Starting the ffmpeg process: ")
	FFmpegProcess(&Env)
	finalChecks(&Env, &Channels, &Proc)
	end := time.Since(start)
	log.Printf("Time taken: %.2f sec", end.Seconds())
}

func FFmpegProcess(Env *Types.TasksEnv) {

	cmd := exec.Command("bash", "./transcoder.sh", Env.INPUT_PATH)
	var stdout, stderr bytes.Buffer
	if os.Getenv("FFMPEG_LOG") == "1" {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	}
	fmt.Printf("Executing command: %s %s %s\n", "bash", "./transcoder.sh", Env.INPUT_PATH)
	fmt.Printf("Working directory: %s\n", getWorkingDir())

	err := cmd.Run()
	if err != nil {
		log.Printf("Command failed with error: %v", err)
		log.Printf("Exit code: %d", cmd.ProcessState.ExitCode())
		//check shell script exists?
		if _, statErr := os.Stat("./transcoder.sh"); os.IsNotExist(statErr) {
			log.Printf("ERROR: transcoder.sh file not found in current directory")
		} else {
			info, _ := os.Stat("./transcoder.sh")
			log.Printf("transcoder.sh exists, permissions: %s", info.Mode())
		}
		log.Fatal(err)
	}
}

func NewEnvs(Env *Types.TasksEnv) {
	godotenv.Load()
	//pass env["MODE"] = test (for tests - uses mock-bucket) [default = hls-bucket]
	mode := os.Getenv("MODE")
	if mode == "test" {
		Env.HLS_BUCKET = os.Getenv("MOCK_BUCKETNAME")
	} else {
		Env.HLS_BUCKET = os.Getenv("HLS_BUCKETNAME")
	}
	Env.TMPFS_PATH = os.Getenv("TMPFS_PATH")
	Env.LOGS_PATH = os.Getenv("LOGS_PATH")
	Env.FILE_ID = os.Getenv("FILE_ID")
	Env.INPUT_PATH = os.Getenv("INPUT_PATH")
	Env.OUT_PATH = os.Getenv("OUT_PATH")
	Env.PROJECT_ID = os.Getenv("PROJECT_ID")
	Env.PUB_TOPIC = os.Getenv("PUB_TOPIC")

	fmt.Printf("ENV gathered: \n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		Env.FILE_ID, Env.INPUT_PATH, Env.OUT_PATH, Env.TMPFS_PATH, Env.HLS_BUCKET, Env.PROJECT_ID, Env.PUB_TOPIC,
	)
}

func NewProcessor(Proc *Types.Processor, Env *Types.TasksEnv) {
	var err error
	Proc.Ctx = context.Background()

	// for gcp: check SA credentials
	if metadata.OnGCE() {
		log.Println("Running on GCP detected via metadata server")
		email, err := metadata.Email("default")
		if err != nil {
			log.Printf("Error retrieving service account email: %v", err)
		} else {
			log.Printf("Effective service account email: %s", email)
		}
	} else {
		log.Printf("Not running on GCP (or metadata unavailable)")
	}

	Proc.Cli, err = storage.NewClient(Proc.Ctx)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	//bucket exists?
	_, err = Proc.Cli.Bucket(Env.HLS_BUCKET).Attrs(Proc.Ctx)
	if err != nil {
		log.Printf("WARNING: Cannot access bucket '%s': %v", Env.HLS_BUCKET, err)
	} else {
		log.Printf("Bucket '%s' exists and is accessible", Env.HLS_BUCKET)
	}

	Proc.Bkt = Proc.Cli.Bucket(Env.HLS_BUCKET)
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
	fmt.Println("Playlists Uploaded")
	fmt.Println("\nPublishing Completion message (on pub-sub topic):")
	GCS.PublishStatus(Env, Env.FILE_ID)
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
	}
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
