package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var verbose = false

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

	gitGetHost := os.Getenv("GIT_GET_HOST")
	if gitGetHost == "" {
		gitGetHost = "github.com"
	}

	verbose = os.Getenv("GIT_GET_VERBOSE") != ""

	gitGetBase := os.Getenv("GIT_GET_BASE")
	if gitGetBase == "" {
		gitGetBase = filepath.Join(os.Getenv("HOME"), "src")
	}

	if verbose {
		log.Println("Using base:", gitGetBase)
	}

	var repoParts []string
	/// if repo.count(":") == 1:
	if strings.Count(repo, ":") == 1 {
		/// if "@" in repo:
		if strings.Contains(repo, "@") {
			///	repo_parts = repo[repo.find("@")+1:].replace(":", "/").split("/")
			repoParts = strings.Split(strings.Replace(repo[strings.Index(repo, "@")+1:len(repo)-1], ":", "/", -1), "/")
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
		repoParts = []string{gitGetHost}
		repoParts = append(repoParts, strings.Split(repo, "/")...)
		/// repo = "git@github.com:%s.git" % repo
		repo = fmt.Sprintf("git@%s:%s.git", gitGetHost, repo)
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
	targetDir = filepath.Join(gitGetBase, targetDir)

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
