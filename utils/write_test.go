package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rng70/versions/v2/canonicalized"
)

// ─── WriteToFile ──────────────────────────────────────────────────────────────

func TestWriteToFile_SimpleStruct(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	data := map[string]string{"key": "value"}
	if err := WriteToFile(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("got %q, want %q", got["key"], "value")
	}
}

func TestWriteToFile_Slice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "slice.json")

	data := []int{1, 2, 3}
	if err := WriteToFile(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var got []int
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(got) != 3 || got[0] != 1 {
		t.Errorf("unexpected content: %v", got)
	}
}

func TestWriteToFile_Nil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nil.json")

	if err := WriteToFile(path, nil); err != nil {
		t.Fatalf("unexpected error writing nil: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	// null is valid JSON
	if string(raw) != "null\n" {
		t.Errorf("nil output: got %q, want %q", string(raw), "null\n")
	}
}

func TestWriteToFile_CreatesSubdirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "out.json")

	if err := WriteToFile(path, "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file was not created")
	}
}

func TestWriteToFile_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overwrite.json")

	// write first time
	if err := WriteToFile(path, 1); err != nil {
		t.Fatalf("first write: %v", err)
	}
	// overwrite with different content
	if err := WriteToFile(path, 99); err != nil {
		t.Fatalf("second write: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var got int
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got != 99 {
		t.Errorf("overwrite: got %d, want 99", got)
	}
}

func TestWriteToFile_VersionSlice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "versions.json")

	v1 := canonicalized.ParseVersionString("1.0.0")
	v2 := canonicalized.ParseVersionString("2.0.0")
	data := []*canonicalized.Version{&v1, &v2}

	if err := WriteToFile(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if len(raw) == 0 {
		t.Error("output file should not be empty")
	}
}

func TestWriteToFile_InvalidPath(t *testing.T) {
	// Use an impossible path (null byte in filename, which is invalid on all OS)
	err := WriteToFile("\x00invalid", "data")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

// ─── WriteToFileWithMinimalContext ────────────────────────────────────────────

func TestWriteToFileWithMinimalContext_WritesOriginalStrings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.json")

	v1 := canonicalized.ParseVersionString("1.0.0")
	v2 := canonicalized.ParseVersionString("2.0.0-beta1")
	data := []*canonicalized.Version{&v1, &v2}

	if err := WriteToFileWithMinimalContext(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var got []string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\ncontent: %s", err, raw)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 strings, got %d: %v", len(got), got)
	}
	if got[0] != "1.0.0" {
		t.Errorf("got[0]: got %q, want %q", got[0], "1.0.0")
	}
	if got[1] != "2.0.0-beta1" {
		t.Errorf("got[1]: got %q, want %q", got[1], "2.0.0-beta1")
	}
}

func TestWriteToFileWithMinimalContext_Empty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")

	if err := WriteToFileWithMinimalContext(path, []*canonicalized.Version{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var got []string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty array, got %v", got)
	}
}

func TestWriteToFileWithMinimalContext_CreatesSubdirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "minimal.json")

	v := canonicalized.ParseVersionString("1.0.0")
	if err := WriteToFileWithMinimalContext(path, []*canonicalized.Version{&v}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("file was not created in subdirectory")
	}
}

func TestWriteToFileWithMinimalContext_PreservesOriginal(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "original.json")

	// v-prefixed version — Original should be "v1.0.0"
	v := canonicalized.ParseVersionString("v1.0.0")
	data := []*canonicalized.Version{&v}

	if err := WriteToFileWithMinimalContext(path, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var got []string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 1 || got[0] != "v1.0.0" {
		t.Errorf("original not preserved: got %v", got)
	}
}

func TestWriteToFileWithMinimalContext_InvalidPath(t *testing.T) {
	err := WriteToFileWithMinimalContext("\x00invalid", []*canonicalized.Version{})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
