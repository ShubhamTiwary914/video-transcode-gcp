package processjob

import (
	"fmt"
	"os"
)

type LogWriter struct {
	File *os.File
	Ch   chan string
	Path string
}

func InitLoggers(loggers []*LogWriter, streams int, channelBufferSize int8) {
	for ctr := 0; ctr < streams; ctr++ {
		dir := fmt.Sprintf("./logs/stream_%d.txt", ctr)
		f, _ := os.Create(dir)
		ch := make(chan string, channelBufferSize)
		loggers[ctr] = &LogWriter{
			File: f,
			Ch:   ch,
			Path: dir,
		}

		//logs receiver go-routines
		go func(lw *LogWriter) {
			for msg := range lw.Ch {
				lw.File.WriteString(msg + "\n")
			}
		}(loggers[ctr])
	}
}

func SetupDirs(streams int, ramfs_path string, logs_path string) {
	for i := 0; i < streams; i++ {
		stream := fmt.Sprintf("stream_%d", i)
		os.MkdirAll(ramfs_path+stream, 0755)
	}
	//reset logs
	os.RemoveAll(logs_path)
	os.MkdirAll(logs_path, 0755)
}
