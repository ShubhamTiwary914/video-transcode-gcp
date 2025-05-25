package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	Types "processjob/types"
	Utils "processjob/utils"

	"cloud.google.com/go/storage"
)

func TestNewEnvs(t *testing.T) {
	env := Types.TasksEnv{}
	NewEnvs(&env)

	bucketExists(t, env.HLS_BUCKET)

	if env.TMPFS_PATH == "" {
		t.Error("TMPFS_PATH is undefined (empty string)")
	} else if _, err := os.Stat(env.TMPFS_PATH); os.IsNotExist(err) {
		t.Errorf("TMPFS_PATH does not exist: %q", env.TMPFS_PATH)
	}

	if env.LOGS_PATH == "" {
		t.Error("LOGS_PATH is undefined (empty string)")
	} else if _, err := os.Stat(env.LOGS_PATH); os.IsNotExist(err) {
		t.Errorf("LOGS_PATH does not exist: %q", env.LOGS_PATH)
	}

	if env.OUT_PATH == "" {
		t.Error("OUT_PATH is undefined (empty string)")
	} else if _, err := os.Stat(env.OUT_PATH); os.IsNotExist(err) {
		t.Errorf("OUT_PATH does not exist: %q", env.OUT_PATH)
	}

	if env.INPUT_PATH == "" {
		t.Error("INPUT_PATH is undefined (empty string)")
	} else if _, err := os.Stat(env.INPUT_PATH); os.IsNotExist(err) {
		t.Errorf("INPUT_PATH does not point to an existing file: %q", env.INPUT_PATH)
	}
}

func TestNewProcessor(t *testing.T) {
	env := Types.TasksEnv{}
	NewEnvs(&env)
	env.HLS_BUCKET = os.Getenv("MOCK_BUCKETNAME")

	if err := Utils.SetupDirs(3, &env); err != nil {
		t.Fatalf("SetupDirs failed: %v", err)
	}

	proc := Types.Processor{}
	NewProcessor(&proc, &env)

	// Check TMPFS stream directories
	for i := 0; i < streams; i++ {
		expected := filepath.Join(env.TMPFS_PATH, fmt.Sprintf("stream_%d", i))
		if stat, err := os.Stat(expected); err != nil || !stat.IsDir() {
			t.Errorf("TMPFS stream dir not created: %s", expected)
		}
	}

	// Check logs/streams dir
	streamLogsPath := filepath.Join(env.LOGS_PATH, "streams")
	if stat, err := os.Stat(streamLogsPath); err != nil || !stat.IsDir() {
		t.Errorf("Stream logs dir not created: %s", streamLogsPath)
	}

	// Check out.log exists
	outLogPath := filepath.Join(env.OUT_PATH, "out.log")
	if _, err := os.Stat(outLogPath); err != nil {
		t.Errorf("Log file not created: %s", outLogPath)
	}

	//check bucket conn
	connBucketName := proc.Bkt.BucketName()
	if connBucketName != env.HLS_BUCKET {
		t.Errorf("Bucket not connected successfully: target: %s, connected: %s", env.HLS_BUCKET, connBucketName)
	}
}

func TestFSWatchers(t *testing.T) {

}

func TestBucket_SingleMockUpload(t *testing.T) {

}

//region Utils

func bucketExists(t *testing.T, bucketName string) error {
	t.Helper()
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create GCS client: %v", err)
	}
	defer client.Close()

	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Errorf("GCS bucket %q does not exist or is inaccessible: %v", bucketName, err)
	}
	return err
}
