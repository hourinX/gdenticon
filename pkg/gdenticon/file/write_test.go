package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func writeText(text string) WriteFunc {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, text)
		return err
	}
}

func TestCreateReplaceAndCreateTemp(t *testing.T) {
	dir := t.TempDir()
	path := Path(dir, "icon", ".txt")
	if err := Create(path, 0o600, writeText("first")); err != nil {
		t.Fatal(err)
	}
	if err := Create(path, 0o600, writeText("duplicate")); err == nil {
		t.Fatal("Create replaced an existing file")
	}
	if err := Replace(path, 0o600, writeText("second")); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Fatalf("file contains %q", got)
	}

	temp, err := CreateTemp(dir, "icon-*", 0o600, writeText("temporary"))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(temp) != dir {
		t.Fatalf("temporary path %q is outside %q", temp, dir)
	}
}

func TestWriteFailuresCleanUp(t *testing.T) {
	dir := t.TempDir()
	fail := func(io.Writer) error { return errors.New("intentional") }
	path := filepath.Join(dir, "failed")
	if err := Create(path, 0o600, fail); err == nil {
		t.Fatal("Create returned nil")
	}
	if _, err := CreateTemp(dir, "failed-*", 0o600, fail); err == nil {
		t.Fatal("CreateTemp returned nil")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(path) {
		t.Fatalf("CreateTemp left a temporary file: %v", entries)
	}
}

func TestPutCache(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	first, err := PutCache(dir, ".png", 0o600, writeText("same content"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := PutCache(dir, ".png", 0o600, writeText("same content"))
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("cache paths differ: %q != %q", first, second)
	}
	if filepath.Ext(first) != ".png" {
		t.Fatalf("cache extension = %q", filepath.Ext(first))
	}
	if _, err := PutCache(dir, ".png", 0o600, func(io.Writer) error { return errors.New("intentional") }); err == nil {
		t.Fatal("PutCache returned nil for write failure")
	}
}

func TestPutCacheEntryReportsCreation(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	first, err := PutCacheEntry(dir, ".png", 0o600, writeText("tracked content"))
	if err != nil {
		t.Fatal(err)
	}
	if !first.Created {
		t.Fatal("first cache write was not reported as created")
	}
	second, err := PutCacheEntry(dir, ".png", 0o600, writeText("tracked content"))
	if err != nil {
		t.Fatal(err)
	}
	if second.Created {
		t.Fatal("existing cache entry was reported as created")
	}
	if first.Path != second.Path {
		t.Fatalf("cache paths differ: %q != %q", first.Path, second.Path)
	}
}

func TestPutCacheEntryRejectsInvalidDirectory(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "not-a-directory")
	if err := os.WriteFile(file, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := PutCacheEntry(filepath.Join(file, "cache"), ".png", 0o600, writeText("content")); err == nil {
		t.Fatal("PutCacheEntry returned nil for invalid directory")
	}
}

func TestMoveNewRejectsMissingSource(t *testing.T) {
	dir := t.TempDir()
	if err := moveNew(filepath.Join(dir, "missing"), filepath.Join(dir, "target"), 0o600); err == nil {
		t.Fatal("moveNew returned nil for missing source")
	}
}

func TestPutCacheEntryConcurrentCreation(t *testing.T) {
	const writers = 8
	dir := filepath.Join(t.TempDir(), "cache")
	results := make(chan CacheEntry, writers)
	errors := make(chan error, writers)

	var wg sync.WaitGroup
	for range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			entry, err := PutCacheEntry(dir, ".png", 0o600, writeText("shared content"))
			if err != nil {
				errors <- err
				return
			}
			results <- entry
		}()
	}
	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		t.Fatal(err)
	}
	created := 0
	path := ""
	for entry := range results {
		if entry.Created {
			created++
		}
		if path == "" {
			path = entry.Path
		} else if entry.Path != path {
			t.Fatalf("cache paths differ: %q != %q", entry.Path, path)
		}
	}
	if created != 1 {
		t.Fatalf("created entries = %d, want 1", created)
	}
}
