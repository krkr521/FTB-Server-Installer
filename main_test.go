package main

import (
	"reflect"
	"testing"

	"ftb-server-downloader/structs"
)

func TestOrderedDownloadSourcesPrefersMirrors(t *testing.T) {
	file := structs.File{
		Url:     "https://primary.invalid/example.jar",
		Mirrors: []string{"https://mirror-one.invalid/example.jar", "https://mirror-two.invalid/example.jar"},
	}
	want := []string{file.Mirrors[0], file.Mirrors[1], file.Url}

	if got := orderedDownloadSources(file); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected source order: got %v, want %v", got, want)
	}
}

func TestDownloadAttemptsCyclesAcrossSources(t *testing.T) {
	file := structs.File{
		Url:     "https://primary.invalid/example.jar",
		Mirrors: []string{"https://mirror.invalid/example.jar"},
	}
	want := []string{file.Mirrors[0], file.Url, file.Mirrors[0], file.Url, file.Mirrors[0], file.Url}

	if got := downloadAttempts(file, 3); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected attempt order: got %v, want %v", got, want)
	}
}
