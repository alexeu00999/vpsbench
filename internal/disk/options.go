package disk

// DiskTarget describes a disk to benchmark.
type DiskTarget struct {
	Device string // device name, e.g. "nvme0n1", "sda"
	Type   string // "NVMe", "SSD", "HDD"
	Path   string // directory path for temp files (mount point or subdir)
}

// Option configures a DiskBench instance.
type Option func(*DiskBench)

// WithTargets sets specific disks to benchmark.
// Each target should have a Path where temp files can be written.
func WithTargets(targets []DiskTarget) Option {
	return func(d *DiskBench) {
		d.targets = targets
	}
}
