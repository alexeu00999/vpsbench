package ram

import (
	"context"
	"testing"
	"time"
)

func TestWorkloadWrite(t *testing.T) {
	ctx := context.Background()
	mbps := workloadWrite(ctx, 500*time.Millisecond)
	if mbps <= 0 {
		t.Errorf("workloadWrite() returned %f, want > 0", mbps)
	}
	t.Logf("Write: %.0f MB/s", mbps)
}

func TestWorkloadRead(t *testing.T) {
	ctx := context.Background()
	mbps := workloadRead(ctx, 500*time.Millisecond)
	if mbps <= 0 {
		t.Errorf("workloadRead() returned %f, want > 0", mbps)
	}
	t.Logf("Read: %.0f MB/s", mbps)
}

func TestWorkloadWriteWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // сразу отменяем

	mbps := workloadWrite(ctx, 5*time.Second)
	// Должен завершиться быстро, без паники
	t.Logf("Write with cancelled ctx: %.0f MB/s", mbps)
}

func TestWorkloadReadWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mbps := workloadRead(ctx, 5*time.Second)
	t.Logf("Read with cancelled ctx: %.0f MB/s", mbps)
}
