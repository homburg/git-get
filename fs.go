package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func init() {
	NoPathError = errors.New("Could not get home path from env vars HOME or USERPROFILE")
}

func defaultPath() string {
	home, err := homePath()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(home, "src")
}

// NoPathError thrown when home path could not automatically be determined
var NoPathError error

func homePath() (string, error) {
	value := ""
	for _, key := range []string{"HOME", "USERPROFILE"} {
		value = os.Getenv(key)
		if value != "" {
			return value, nil
		}
	}

	return "", NoPathError
}

func dirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func rmDir(path string) error {
	if verbose {
		log.Printf("rmDir?: %s\n", path)
	}

	if ok, _ := dirIsEmpty(path); ok {
		if verbose {
			log.Printf("Removing %s\n", path)
		}
		return os.Remove(path)
	}

	return nil
}
