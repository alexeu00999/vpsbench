package network

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMeasurePing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx := context.Background()
	latency, err := measurePing(ctx, ts.URL)
	if err != nil {
		t.Fatalf("measurePing failed: %v", err)
	}

	if latency <= 0 {
		t.Errorf("expected positive latency, got %f", latency)
	}
}

func TestMeasureDownload(t *testing.T) {
	// Создаем 1MB данных для теста
	data := make([]byte, 1024*1024)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer ts.Close()

	ctx := context.Background()
	speed, err := measureDownload(ctx, ts.URL)
	if err != nil {
		t.Fatalf("measureDownload failed: %v", err)
	}

	if speed <= 0 {
		t.Errorf("expected positive speed, got %f", speed)
	}
}

func TestMeasureUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx := context.Background()
	speed, err := measureUpload(ctx, ts.URL)
	if err != nil {
		t.Fatalf("measureUpload failed: %v", err)
	}

	if speed <= 0 {
		t.Errorf("expected positive speed, got %f", speed)
	}
}

func TestNetworkBench_Run(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	srvs := []Server{
		{ID: "test-srv", Region: "TEST", Provider: "Local", URL: ts.URL},
	}

	nb := &NetworkBench{servers: srvs}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := nb.Run(ctx)
	if res.Err != nil {
		t.Fatalf("benchmark failed: %v", res.Err)
	}

	if len(res.Results) == 0 {
		t.Errorf("expected results, got 0")
	}

	foundPing := false
	for _, r := range res.Results {
		if r.Name == "Ping (TEST)" {
			foundPing = true
			break
		}
	}

	if !foundPing {
		t.Errorf("Ping (TEST) result not found")
	}
}
