package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rng70/versions/canonicalized"
)

func generateJsonEncoder(filename string) (*json.Encoder, *os.File, error) {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return nil, nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil, err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")

	return enc, file, nil
}

func WriteToFile(filename string, out any) error {
	enc, file, err := generateJsonEncoder(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := enc.Encode(out); err != nil {
		return err
	}

	return nil
}

func WriteToFileWithMinimalContext(filename string, out []*canonicalized.Version) error {
	minimalEnc, file, err := generateJsonEncoder(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// minimal context: only Original and Parsed.Canonical
	minimalOut := make([]string, len(out))
	for i, v := range out {
		minimalOut[i] = v.Original
	}
	if err := minimalEnc.Encode(minimalOut); err != nil {
		return err
	}

	return nil
}
