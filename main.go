package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var verbose = false

// NoPathError thrown when home path could not automatically be determined
var NoPathError error

func init() {
	NoPathError = errors.New("Could not get home path from env vars HOME or USERPROFILE")
}

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

func runGit(args ...string) {
	fullArgs := []string{"git"}
	fullArgs = append(fullArgs, args...)
	if verbose {
		log.Println("Running:", fullArgs)
	}
	syscall.Exec("/usr/bin/git", fullArgs, os.Environ())
}

func main() {

	// if len(sys.argv) <= 1:
	// 	os.execlp("git", "git")

	if len(os.Args) <= 1 {
		runGit()
	}

	/// repo = str(sys.argv[-1])
	repo := os.Args[len(os.Args)-1]

	/// if repo.endswith(".git"):
	/// 	repo = repo[:-4]
	repo = strings.TrimSuffix(repo, ".git")

	/// if "/" not in repo:
	/// 	os.execlp("git", "git")
	if !strings.Contains(repo, "/") {
		runGit()
	}

	host := os.Getenv("GIT_GET_HOST")
	if host == "" {
		host = "github.com"
	}

	verbose = os.Getenv("GIT_GET_VERBOSE") != ""

	path := os.Getenv("GIT_GET_PATH")
	if path == "" {
		home, err := homePath()
		if err != nil {
			log.Fatal(err)
		}
		path = filepath.Join(home, "src")
	}

	if verbose {
		log.Println("Using base:", path)
	}

	var repoParts []string
	/// if repo.count(":") == 1:
	if strings.Count(repo, ":") == 1 {
		/// if "@" in repo:
		if strings.Contains(repo, "@") {
			///	repo_parts = repo[repo.find("@")+1:].replace(":", "/").split("/")
			repoParts = strings.Split(strings.Replace(repo[strings.Index(repo, "@")+1:len(repo)], ":", "/", -1), "/")
		} else {
			///	repo_parts = repo.replace(":", "/").split("/")
			repoParts = strings.Split(strings.Replace(repo, ":", "/", -1), "/")
			///	repo = "git@" + repo
			repo = "git@" + repo
		}
		/// elif repo.count("/") == 1:
	} else if strings.Count(repo, "/") == 1 {
		///	# Something from github
		///	repo_parts = ["github.com"] + repo.split("/")
		repoParts = []string{host}
		repoParts = append(repoParts, strings.Split(repo, "/")...)
		/// repo = "git@github.com:%s.git" % repo
		repo = fmt.Sprintf("git@%s:%s.git", host, repo)
		if verbose {
			log.Println("Did build repo:", repo)
		}
	}

	/// else:
	/// 	# http repo?
	/// 	print("TODO...")
	/// 	exit(1)
	///

	targetDir := filepath.Join(repoParts...)
	targetDir = filepath.Join(path, targetDir)

	if verbose {
		log.Println("Using target dir:", targetDir)
	}
	err := os.MkdirAll(targetDir, os.ModePerm)
	if nil != err {
		log.Fatal(err)
	}
	/// target_dir = os.path.join(os.getenv("HOME"), "src", *repo_parts)
	/// try:
	/// 	os.makedirs(target_dir)
	/// except OSError as exc:
	/// 	if exc.errno == errno.EEXIST and os.path.isdir(target_dir):
	/// 		pass
	/// 	else:
	/// 		raise
	///
	/// os.execlp("git", "git", "clone", repo, target_dir)
	runGit("clone", repo, targetDir)
}
