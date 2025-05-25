package processjob

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	Types "processjob/types"
	Utils "processjob/utils"
)

// offloads ffmpeg -> (tmpfs/tmpfs) -> GCS bucket
func GCS_offloader(Env *Types.TasksEnv, Channels *Types.ChannelsContainer, Proc *Types.Processor, index int) {
	for {
		select {
		case event, ok := <-Proc.Watchers[index].Events:
			if !ok {
				return
			}
			Channels.Loggers[index].Ch <- string(event.Name + " " + event.Op.String())
			//start upload after next CREATE
			if event.Op.String() == "CREATE" {
				dir, _, ext := Utils.GetFilePath_Split(event.Name)
				targetFile := fmt.Sprintf("%s%04d.ts", dir, Proc.ProcessedCtr[index]-1)
				if ext == "m3u8" || Proc.ProcessedCtr[index] > 0 {
					Channels.UploadCh <- Types.UploadEvent{FilePath: targetFile, StreamID: index, FileID: Env.FILE_ID}
				}
				if ext != "m3u8" {
					Proc.ProcessedCtr[index]++
				}
			}
		case err, ok := <-Proc.Watchers[index].Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

// upload files to HLS bucket
func GCS_uploader(Proc *Types.Processor, Env *Types.TasksEnv, localFile string, streamID int) {
	start := time.Now()
	data, err := os.ReadFile(localFile)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}
	//filepath as it will appear in GCS bucket
	filebase := filepath.Base(localFile)
	gsFile := fmt.Sprintf("%s/%s/%s", Env.FILE_ID, Utils.StreamResolutions[streamID], filebase)
	//root of tmpfs (for master playlist)
	if streamID == -1 {
		gsFile = fmt.Sprintf("%s/%s", Env.FILE_ID, filebase)
	}
	//write to HLS-bucket
	obj := Proc.Bkt.Object(gsFile)
	w := obj.NewWriter(Proc.Ctx)
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
