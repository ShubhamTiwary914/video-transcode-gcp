package processjob

import (
	"fmt"
	"os"
	"path/filepath"
	Types "processjob/types"
	"strconv"
	"strings"
)

// higher->lower
var StreamResolutions = map[int]string{
	0: "stream_0", //1080p
	1: "stream_1", //720p
	2: "stream_2", //480p
}

func InitLoggers(loggers []*Types.LogWriter, streams int, logs_path string, channelBufferSize int8) {
	for ctr := 0; ctr < streams; ctr++ {
		dir := fmt.Sprintf("%s/streams/stream_%d.txt", logs_path, ctr)
		f, _ := os.Create(dir)
		ch := make(chan string, channelBufferSize)
		loggers[ctr] = &Types.LogWriter{
			File: f,
			Ch:   ch,
			Path: dir,
		}
	}
}

func SetupDirs(streams int, Env *Types.TasksEnv) error {
	//clean state
	os.RemoveAll(Env.TMPFS_PATH)
	os.RemoveAll(Env.LOGS_PATH)
	//tmpfs dir: stream_i
	for i := 0; i < streams; i++ {
		dir := filepath.Join(Env.TMPFS_PATH, fmt.Sprintf("stream_%d", i))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}
	//stream logs
	streamsLog := filepath.Join(Env.LOGS_PATH, "streams")
	if err := os.MkdirAll(streamsLog, 0755); err != nil {
		return fmt.Errorf("creating %s: %w", streamsLog, err)
	}
	//stdout/out.log
	if err := os.MkdirAll(Env.OUT_PATH, 0755); err != nil {
		return fmt.Errorf("creating %s: %w", Env.OUT_PATH, err)
	}
	outLog := filepath.Join(Env.OUT_PATH, "out.log")
	f, err := os.Create(outLog)
	if err != nil {
		return fmt.Errorf("creating %s: %w", outLog, err)
	}
	return f.Close()
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
