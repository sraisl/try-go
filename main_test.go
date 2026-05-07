package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFilterDirsMatchesAllTerms(t *testing.T) {
	dirs := []string{
		"2026-04-10-csharp",
		"2026-04-10-gameboy",
		"2026-05-05-go",
	}

	got := filterDirs(dirs, "04 game")
	want := []string{"2026-04-10-gameboy"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterDirs() = %v, want %v", got, want)
	}
}

func TestFilterDirsIsCaseInsensitive(t *testing.T) {
	dirs := []string{
		"2026-04-10-CSharp",
		"2026-04-10-gameboy",
		"2026-05-05-go",
	}

	got := filterDirs(dirs, "csharp")
	want := []string{"2026-04-10-CSharp"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterDirs() = %v, want %v", got, want)
	}
}

func TestFilterDirsReturnsAllForBlankQuery(t *testing.T) {
	dirs := []string{"one", "two"}

	got := filterDirs(dirs, "   ")

	if !reflect.DeepEqual(got, dirs) {
		t.Fatalf("filterDirs() = %v, want %v", got, dirs)
	}
}

func TestApplyFilterResetsCursor(t *testing.T) {
	m := model{
		allDirs: []string{"2026-04-10-gameboy", "2026-05-05-go"},
		query:   "go",
		cursor:  1,
	}

	m.applyFilter()

	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", m.cursor)
	}

	want := []string{"2026-05-05-go"}
	if !reflect.DeepEqual(m.dirs, want) {
		t.Fatalf("dirs = %v, want %v", m.dirs, want)
	}
}

func TestSanitizeDirName(t *testing.T) {
	got := sanitizeDirName("  Hello, Go World!  ")
	want := "hello-go-world"

	if got != want {
		t.Fatalf("sanitizeDirName() = %q, want %q", got, want)
	}
}

func TestSanitizeDirNameKeepsLettersAndDigits(t *testing.T) {
	got := sanitizeDirName("Äpfel 123 & Go")
	want := "äpfel-123-go"

	if got != want {
		t.Fatalf("sanitizeDirName() = %q, want %q", got, want)
	}
}

func TestSanitizeDirNameReturnsEmptyForInvalidInput(t *testing.T) {
	if got := sanitizeDirName("!!!"); got != "" {
		t.Fatalf("sanitizeDirName() = %q, want empty string", got)
	}
}

func TestRemoveDirName(t *testing.T) {
	dirs := []string{"one", "two", "three"}

	got := removeDirName(dirs, "two")
	want := []string{"one", "three"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("removeDirName() = %v, want %v", got, want)
	}
}

func TestRemoveDirDeletesDirectory(t *testing.T) {
	root := t.TempDir()
	name := "2026-05-07-test"
	path := filepath.Join(root, name)

	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	if err := os.WriteFile(filepath.Join(path, "note.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := removeDir(root, name); err != nil {
		t.Fatalf("removeDir() error = %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Stat() error = %v, want not exist", err)
	}
}

func TestRemoveDirRejectsEscapingRoot(t *testing.T) {
	root := t.TempDir()

	if err := removeDir(root, "../outside"); err == nil {
		t.Fatal("removeDir() error = nil, want error")
	}
}
