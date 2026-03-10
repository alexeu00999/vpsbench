package disk

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSequentialWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dir := t.TempDir()
	speed, err := testSequentialWrite(ctx, dir)
	if err != nil {
		t.Fatalf("testSequentialWrite() error: %v", err)
	}
	if speed <= 0 {
		t.Errorf("testSequentialWrite() = %f, want > 0", speed)
	}
	t.Logf("Sequential write speed: %.2f MB/s", speed)

	// Verify temp files are cleaned up
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vpsbench-disk-w-") {
			t.Errorf("temp file not cleaned up: %s", filepath.Join(dir, e.Name()))
		}
	}
}

func TestSequentialRead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dir := t.TempDir()
	speed, err := testSequentialRead(ctx, dir)
	if err != nil {
		t.Fatalf("testSequentialRead() error: %v", err)
	}
	if speed <= 0 {
		t.Errorf("testSequentialRead() = %f, want > 0", speed)
	}
	t.Logf("Sequential read speed: %.2f MB/s", speed)

	// Verify temp files are cleaned up
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vpsbench-disk-r-") {
			t.Errorf("temp file not cleaned up: %s", filepath.Join(dir, e.Name()))
		}
	}
}

func TestRandom4KIOPS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dir := t.TempDir()
	iops, err := testRandom4KIOPS(ctx, dir)
	if err != nil {
		t.Fatalf("testRandom4KIOPS() error: %v", err)
	}
	if iops <= 0 {
		t.Errorf("testRandom4KIOPS() = %f, want > 0", iops)
	}
	t.Logf("Random 4K IOPS: %.2f", iops)

	// Verify temp files are cleaned up
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vpsbench-disk-4k-") {
			t.Errorf("temp file not cleaned up: %s", filepath.Join(dir, e.Name()))
		}
	}
}

func TestSequentialWriteContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dir := t.TempDir()
	// Should not hang on cancelled context
	_, _ = testSequentialWrite(ctx, dir)
}

func TestSequentialReadContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dir := t.TempDir()
	_, _ = testSequentialRead(ctx, dir)
}
