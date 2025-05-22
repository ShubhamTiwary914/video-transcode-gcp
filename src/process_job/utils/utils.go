package processjob

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

var StreamResolutions = map[int]string{
	0: "1080p",
	1: "720p",
	2: "480p",
}

func InitLoggers(loggers []*LogWriter, streams int, logs_path string, channelBufferSize int8) {
	for ctr := 0; ctr < streams; ctr++ {
		dir := fmt.Sprintf("%s/streams/stream_%d.txt", logs_path, ctr)
		f, _ := os.Create(dir)
		ch := make(chan string, channelBufferSize)
		loggers[ctr] = &LogWriter{
			File: f,
			Ch:   ch,
			Path: dir,
		}
	}
}

func SetupDirs(streams int, tmpfs_path string, logs_path string) {
	for i := 0; i < streams; i++ {
		os.MkdirAll(fmt.Sprintf("%s/stream_%d", tmpfs_path, i), 0755)
	}
	//reset stream logs
	streamlogsPath := fmt.Sprintf("%s/streams", logs_path)
	os.RemoveAll(streamlogsPath)
	os.MkdirAll(streamlogsPath, 0755)
}

func GetFilePath_Split(path string) (string, string, string) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return dir + string(os.PathSeparator), name, strings.TrimPrefix(ext, ".")
}

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}
