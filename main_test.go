package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDirIsEmpty(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
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
		t.Errorf(`Expected %s to be empty, was not.`, path)
	}

	parent := filepath.Dir(path)
	ok, err = dirIsEmpty(parent)

	if nil != err {
		t.Fatal(err)
	}

	if ok {
		t.Errorf(`Expected %s not to be empty, was.`, parent)
	}
}

func TestRepoHasScheme(t *testing.T) {
	testData := map[string]bool{
		"git@github.com:some/thing.git":              false,
		"http://what.dk/ever.git":                    true,
		"https://github.com/something/different.git": true,
	}

	for addr, expected := range testData {
		if hasScheme(addr) != expected {
			t.Errorf(`Expected scheme presence: %t for "%s", got %t.`, expected, addr, !expected)
		}
	}
}

type parseRepoPair struct {
	repo      string
	repoParts []string
}

func TestParseRepo(t *testing.T) {
	testData := map[string]parseRepoPair{
		"homburg/tree": {
			"git@gitgit.git:homburg/tree.git",
			[]string{"gitgit.git", "homburg", "tree"},
		},
		"github.com:homburg/branch": {
			"git@github.com:homburg/branch.git",
			[]string{"github.com", "homburg", "branch"},
		},
		"github.com/homburg/stick": {
			"git@gitgit.git:github.com/homburg/stick.git",
			[]string{"gitgit.git", "github.com", "homburg", "stick"},
		},
		"https://github.com/some/thing.git": {
			"https://github.com/some/thing",
			[]string{"github.com", "some", "thing"},
		},
		"http://some.where/test/it": {
			"http://some.where/test/it",
			[]string{"some.where", "test", "it"},
		},
	}

	for repo, result := range testData {
		_, parsedRepo, repoParts := repoParsers.parse(repo, "gitgit.git")

		if !reflect.DeepEqual(result, parseRepoPair{parsedRepo, repoParts}) {
			t.Errorf(`Expected repo: %s, and repo parts: %q, got repo: %s and repo parts: %q.`, result.repo, result.repoParts, parsedRepo, repoParts)
		}
	}
}
