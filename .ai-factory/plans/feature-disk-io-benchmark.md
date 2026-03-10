# Implementation Plan: Disk I/O Benchmark

Branch: feature/disk-io-benchmark
Created: 2026-03-10

## Settings
- Testing: yes
- Logging: verbose
- Docs: yes

## Roadmap Linkage
Milestone: "Disk I/O бенчмарк — sequential read/write, random 4K IOPS, автодетект дисков"
Rationale: This is the next unchecked milestone in the roadmap, directly following the completed RAM benchmark.

## Commit Plan
- **Commit 1** (after tasks 1-2): "feat(disk): refactor workloads and add random 4K IOPS test"
- **Commit 2** (after tasks 3-5): "feat(disk): add per-disk benchmarking with platform I/O optimizations"
- **Commit 3** (after task 6): "test(disk): add unit tests for disk benchmark module"

## Tasks

### Phase 1: Refactor & Core Implementation

- [x] Task 1: Extract disk workloads into workloads.go
  - Move `testSequentialWrite` and `testSequentialRead` from disk.go to new `internal/disk/workloads.go`
  - Increase test file size from 64 MB to 256 MB for accuracy
  - Add fsync after write, drop OS cache before read
  - Follow pattern from cpu/workloads.go and ram/workloads.go
  - Files: `internal/disk/workloads.go` (new), `internal/disk/disk.go`

- [x] Task 2: Implement random 4K IOPS test (depends on 1)
  - Add `testRandom4KIOPS()` to `internal/disk/workloads.go`
  - Pre-allocate 128 MB test file, random 4K read/write at random offsets for 5s each
  - Use O_SYNC or fsync per write for real IOPS measurement
  - Add `bench.Result{Name: "Rand 4K IOPS", Value: iops, Unit: "IOPS"}` to Run()
  - Matches existing baseline key `"DISK:Rand 4K IOPS": 1000000` in rating.go
  - Files: `internal/disk/workloads.go`, `internal/disk/disk.go`

<!-- Commit checkpoint: tasks 1-2 -->

### Phase 2: Multi-disk & Platform Optimizations

- [ ] Task 3: Add per-disk benchmarking support (depends on 1, 2)
  - Accept disk info via functional options: `New(opts ...Option)`
  - Test each detected disk on its mount point
  - Aggregate per-disk results: "Seq. Write (nvme0n1)", "Rand 4K IOPS (sda)"
  - Fallback to default temp dir if no disks detected
  - Files: `internal/disk/disk.go`, `internal/disk/workloads.go`, `internal/disk/options.go` (new)

- [ ] Task 4: Add platform-specific I/O optimizations (depends on 1)
  - `internal/disk/io_linux.go`: O_DIRECT, syscall.Fadvise FADV_DONTNEED, 4096-aligned buffers
  - `internal/disk/io_darwin.go`: F_NOCACHE via fcntl, graceful fallback
  - `internal/disk/io_default.go`: no-op implementations for other platforms
  - Files: `internal/disk/io_linux.go`, `internal/disk/io_darwin.go`, `internal/disk/io_default.go` (all new)

- [ ] Task 5: Integrate disk info into DiskBench from main.go (depends on 3)
  - Pass `info.Disks` to `disk.New()` via options
  - Remove post-hoc Info override for DISK module in main.go
  - Let DiskBench set its own Info from tested disks
  - Files: `cmd/vpsbench/main.go`, `internal/disk/disk.go`

<!-- Commit checkpoint: tasks 3-5 -->

### Phase 3: Testing

- [ ] Task 6: Write unit tests for disk benchmark (depends on 2, 3, 4)
  - `internal/disk/disk_test.go`: test Name(), Run(), result names, context cancellation
  - `internal/disk/workloads_test.go`: test each workload returns positive values, temp file cleanup
  - Use `t.TempDir()`, short durations (1s) for CI
  - Files: `internal/disk/disk_test.go`, `internal/disk/workloads_test.go` (both new)

<!-- Commit checkpoint: task 6 -->
