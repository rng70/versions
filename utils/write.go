package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rng70/versions/canonicalized"
)

func WriteToFile(filename string, out []*canonicalized.Version) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	//enc := json.NewEncoder(os.Stdout)
	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")
	if err := enc.Encode(out); err != nil {
		os.Exit(1)
	}

	return err
}

func WriteToFileWithMinimalContext(filename string, out []*canonicalized.Version) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "    ")

	// minimal context: only Original and Parsed.Canonical
	minimalOut := make([]string, len(out))
	for i, v := range out {
		minimalOut[i] = v.Original
	}
	if err := enc.Encode(minimalOut); err != nil {
		os.Exit(1)
	}

	return err
}
