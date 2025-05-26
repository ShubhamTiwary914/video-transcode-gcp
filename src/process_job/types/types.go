package processjob

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
)

type TasksEnv struct {
	HLS_BUCKET string
	TMPFS_PATH string
	LOGS_PATH  string
	FILE_ID    string
	INPUT_PATH string
	OUT_PATH   string
}

type Processor struct {
	Ctx          context.Context
	Bkt          *storage.BucketHandle
	Cli          *storage.Client
	UploadWg     sync.WaitGroup
	Watchers     []*fsnotify.Watcher
	ProcessedCtr map[int]int
}

type ChannelsContainer struct {
	UploadCh chan UploadEvent
	Loggers  []*LogWriter
}

type LogWriter struct {
	File *os.File
	Ch   chan string
	Path string
}

type UploadEvent struct {
	FilePath string
	StreamID int
	FileID   string
}
