package disk

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	mrand "math/rand/v2"
	"os"
	"time"
)

const (
	// seqBlockSize — размер блока для последовательных тестов (1 MB).
	seqBlockSize = 1024 * 1024

	// seqTotalBlocks — количество блоков для последовательных тестов (256 MB).
	seqTotalBlocks = 256

	// rand4KBlockSize — размер блока для случайного I/O (4 KB).
	rand4KBlockSize = 4096

	// rand4KFileSize — размер файла для случайного I/O (128 MB).
	rand4KFileSize = 128 * 1024 * 1024

	// rand4KDuration — длительность каждого теста (чтение/запись).
	rand4KDuration = 5 * time.Second
)

// testSequentialWrite измеряет скорость последовательной записи на диск.
// Записывает seqTotalBlocks блоков по seqBlockSize, затем вызывает fsync.
// Возвращает скорость в MB/s.
func testSequentialWrite(ctx context.Context, dir string) (float64, error) {
	slog.Debug("[disk.seqWrite] starting",
		"block_size", seqBlockSize,
		"total_blocks", seqTotalBlocks,
		"total_mb", seqTotalBlocks,
		"dir", dir,
	)

	data := make([]byte, seqBlockSize)
	if _, err := rand.Read(data); err != nil {
		return 0, fmt.Errorf("generate random data: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "vpsbench-disk-w-*")
	if err != nil {
		return 0, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	defer tmpFile.Close()

	slog.Debug("[disk.seqWrite] temp file created", "path", tmpPath)

	start := time.Now()
	for i := 0; i < seqTotalBlocks; i++ {
		if ctx.Err() != nil {
			slog.Debug("[disk.seqWrite] cancelled by context", "blocks_written", i)
			break
		}
		if _, err := tmpFile.Write(data); err != nil {
			return 0, fmt.Errorf("write block %d: %w", i, err)
		}
	}

	// Fsync to ensure data is flushed to disk
	syncStart := time.Now()
	if err := tmpFile.Sync(); err != nil {
		return 0, fmt.Errorf("fsync: %w", err)
	}
	syncElapsed := time.Since(syncStart)

	elapsed := time.Since(start)
	mbPerSec := float64(seqTotalBlocks) / elapsed.Seconds()

	slog.Debug("[disk.seqWrite] finished",
		"elapsed", elapsed,
		"sync_time", syncElapsed,
		"mb_per_sec", mbPerSec,
	)

	return mbPerSec, nil
}

// testSequentialRead измеряет скорость последовательного чтения с диска.
// Сначала записывает файл, затем закрывает и переоткрывает его для чтения
// (чтобы по возможности сбросить page cache). Возвращает скорость в MB/s.
func testSequentialRead(ctx context.Context, dir string) (float64, error) {
	slog.Debug("[disk.seqRead] starting",
		"block_size", seqBlockSize,
		"total_blocks", seqTotalBlocks,
		"total_mb", seqTotalBlocks,
		"dir", dir,
	)

	// Prepare test file
	data := make([]byte, seqBlockSize)
	if _, err := rand.Read(data); err != nil {
		return 0, fmt.Errorf("generate random data: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "vpsbench-disk-r-*")
	if err != nil {
		return 0, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	slog.Debug("[disk.seqRead] writing test data", "path", tmpPath)

	for i := 0; i < seqTotalBlocks; i++ {
		if _, err := tmpFile.Write(data); err != nil {
			tmpFile.Close()
			return 0, fmt.Errorf("write block %d: %w", i, err)
		}
	}
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return 0, fmt.Errorf("fsync: %w", err)
	}
	tmpFile.Close()

	// Reopen for reading — helps avoid reading from page cache on some OSes
	readFile, err := os.Open(tmpPath)
	if err != nil {
		return 0, fmt.Errorf("reopen for read: %w", err)
	}
	defer readFile.Close()

	// Try to drop page cache (platform-specific, best effort)
	dropPageCache(readFile)

	buf := make([]byte, seqBlockSize)
	start := time.Now()
	for i := 0; i < seqTotalBlocks; i++ {
		if ctx.Err() != nil {
			slog.Debug("[disk.seqRead] cancelled by context", "blocks_read", i)
			break
		}
		if _, err := readFile.Read(buf); err != nil {
			slog.Debug("[disk.seqRead] read ended", "blocks_read", i, "reason", err)
			break
		}
	}
	elapsed := time.Since(start)
	mbPerSec := float64(seqTotalBlocks) / elapsed.Seconds()

	slog.Debug("[disk.seqRead] finished",
		"elapsed", elapsed,
		"mb_per_sec", mbPerSec,
	)

	return mbPerSec, nil
}

// testRandom4KIOPS измеряет случайный 4K I/O (IOPS).
// Создаёт файл rand4KFileSize, затем выполняет случайные 4K записи и чтения
// в течение rand4KDuration каждый. Возвращает среднее IOPS (write+read)/2.
func testRandom4KIOPS(ctx context.Context, dir string) (float64, error) {
	maxOffset := int64(rand4KFileSize - rand4KBlockSize)

	slog.Debug("[disk.rand4k] starting",
		"file_size_mb", rand4KFileSize/1024/1024,
		"block_size", rand4KBlockSize,
		"max_offsets", rand4KFileSize/rand4KBlockSize,
		"duration_per_phase", rand4KDuration,
		"dir", dir,
	)

	// Create and pre-allocate test file
	tmpFile, err := os.CreateTemp(dir, "vpsbench-disk-4k-*")
	if err != nil {
		return 0, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	defer tmpFile.Close()

	// Pre-fill with random data
	slog.Debug("[disk.rand4k] pre-filling test file", "path", tmpPath)
	fillBuf := make([]byte, seqBlockSize) // reuse 1 MB block size
	if _, err := rand.Read(fillBuf); err != nil {
		return 0, fmt.Errorf("generate fill data: %w", err)
	}
	fillBlocks := rand4KFileSize / seqBlockSize
	for i := 0; i < fillBlocks; i++ {
		if _, err := tmpFile.Write(fillBuf); err != nil {
			return 0, fmt.Errorf("pre-fill block %d: %w", i, err)
		}
	}
	if err := tmpFile.Sync(); err != nil {
		return 0, fmt.Errorf("pre-fill fsync: %w", err)
	}

	// Random 4K write IOPS
	slog.Debug("[disk.rand4k] starting write phase")
	writeBuf := make([]byte, rand4KBlockSize)
	if _, err := rand.Read(writeBuf); err != nil {
		return 0, fmt.Errorf("generate write data: %w", err)
	}

	writeOps := 0
	writeDeadline := time.Now().Add(rand4KDuration)
	for time.Now().Before(writeDeadline) {
		if ctx.Err() != nil {
			slog.Debug("[disk.rand4k] write phase cancelled", "ops", writeOps)
			break
		}
		offset := mrand.Int64N(maxOffset)
		// Align offset to 4K boundary
		offset = (offset / rand4KBlockSize) * rand4KBlockSize
		if _, err := tmpFile.WriteAt(writeBuf, offset); err != nil {
			return 0, fmt.Errorf("random write at offset %d: %w", offset, err)
		}
		// Sync each write for real IOPS measurement
		if err := tmpFile.Sync(); err != nil {
			return 0, fmt.Errorf("fsync after random write: %w", err)
		}
		writeOps++
	}
	writeIOPS := float64(writeOps) / rand4KDuration.Seconds()
	slog.Debug("[disk.rand4k] write phase finished", "ops", writeOps, "iops", writeIOPS)

	// Random 4K read IOPS
	slog.Debug("[disk.rand4k] starting read phase")
	readBuf := make([]byte, rand4KBlockSize)
	readOps := 0
	readDeadline := time.Now().Add(rand4KDuration)
	for time.Now().Before(readDeadline) {
		if ctx.Err() != nil {
			slog.Debug("[disk.rand4k] read phase cancelled", "ops", readOps)
			break
		}
		offset := mrand.Int64N(maxOffset)
		offset = (offset / rand4KBlockSize) * rand4KBlockSize
		if _, err := tmpFile.ReadAt(readBuf, offset); err != nil {
			return 0, fmt.Errorf("random read at offset %d: %w", offset, err)
		}
		readOps++
	}
	readIOPS := float64(readOps) / rand4KDuration.Seconds()
	slog.Debug("[disk.rand4k] read phase finished", "ops", readOps, "iops", readIOPS)

	// Average of write and read IOPS
	avgIOPS := (writeIOPS + readIOPS) / 2

	slog.Info("[disk.rand4k] completed",
		"write_iops", writeIOPS,
		"read_iops", readIOPS,
		"avg_iops", avgIOPS,
	)

	return avgIOPS, nil
}
