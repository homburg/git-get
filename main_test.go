package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestDirIsEmpty(t *testing.T) {
	log.Println(os.Args)
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(cwd, "tmp-test/parent-dir/empty-dir")

	err = os.MkdirAll(path, os.ModePerm)
	if nil != err {
		t.Fatal(err)
	}

	ok, err := dirIsEmpty(path)
	if nil != err {
		t.Fatal(err)
	}

	if !ok {
		t.Fatalf(`Expected %s to be empty, was not.`, path)
	}

	parent := filepath.Dir(path)
	ok, err = dirIsEmpty(parent)

	if nil != err {
		t.Fatal(err)
	}

	if ok {
		t.Fatalf(`Expected %s not to be empty, was.`, parent)
	}
}
