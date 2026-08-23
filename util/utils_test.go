package util

import (
	"testing"

	"ftb-server-downloader/structs"
)

func TestParseInstallerName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want1 int
		want2 int
	}{
		{"installer_123 should be 123", "installer_123", 123, 0},
		{"installer_123-1234 should be 123", "installer_123-1234", 123, 0},
		{"installer_123_1234 should be 123, 1234", "installer_123_1234", 123, 1234},
		{"installer_123_1234_5678 should be 123, 1234", "installer_123_1234_5678", 123, 1234},
		{"installer-123-1234 should be error", "installer-123-1234", 0, 0},
		{"installer should be error", "installer", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pack, version, err := ParseInstallerName(tt.input)
			if err != nil && err.Error() != "invalid installer name" && tt.want1 != 0 && tt.want2 != 0 {
				t.Errorf("got unexpected error %s", err)
			}
			if pack != tt.want1 {
				t.Errorf("got %d, want %d", pack, tt.want1)
			}
			if version != tt.want2 {
				t.Errorf("got %d, want %d", version, tt.want2)
			}
		})
	}
}

func TestFailedDownloadHandlerSwitchesImmediately(t *testing.T) {
	file := structs.File{Name: "example.jar"}
	mirrors := []string{"https://primary.invalid/example.jar", "https://mirror.invalid/example.jar"}

	retry, next, err := FailedDownloadHandler(0, 0, file, mirrors[0], mirrors)
	if err != nil {
		t.Fatalf("unexpected error while another mirror is available: %v", err)
	}
	if retry || !next {
		t.Fatalf("expected immediate mirror switch, got retry=%v next=%v", retry, next)
	}
}

func TestFailedDownloadHandlerFailsAfterLastMirror(t *testing.T) {
	file := structs.File{Name: "example.jar"}
	mirrors := []string{"https://primary.invalid/example.jar"}

	retry, next, err := FailedDownloadHandler(0, 0, file, mirrors[0], mirrors)
	if err == nil {
		t.Fatal("expected an error after the last mirror failed")
	}
	if retry || next {
		t.Fatalf("expected no further attempt, got retry=%v next=%v", retry, next)
	}
}
