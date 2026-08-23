package util

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDownloadTimeoutTracksIdleTimeNotTotalDuration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		flusher := w.(http.Flusher)
		for range 4 {
			_, _ = w.Write([]byte("chunk"))
			flusher.Flush()
			time.Sleep(50 * time.Millisecond)
		}
	}))
	defer server.Close()

	dl, err := NewDownload(filepath.Join(t.TempDir(), "active-download.bin"), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	dl.Timeout = 100 * time.Millisecond

	if err := dl.Do(); err != nil {
		t.Fatalf("active download exceeded the timeout but should have succeeded: %v", err)
	}
}

func TestDownloadTimeoutCancelsStalledTransfer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		time.Sleep(250 * time.Millisecond)
		_, _ = w.Write([]byte("late"))
	}))
	defer server.Close()

	dl, err := NewDownload(filepath.Join(t.TempDir(), "stalled-download.bin"), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	dl.Timeout = 100 * time.Millisecond

	err = dl.Do()
	if err == nil || !strings.Contains(err.Error(), "no download progress") {
		t.Fatalf("expected an idle-timeout error, got %v", err)
	}
}
