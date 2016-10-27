package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var verbose = false

// NoPathError thrown when home path could not automatically be determined
var NoPathError error

type exitStatusError struct {
	error
	exitCode int
}

func (err exitStatusError) exit() {
	if err.exitCode != 0 {
		os.Exit(err.exitCode)
	} else {
		log.Fatal(err)
	}
}

func init() {
	NoPathError = errors.New("Could not get home path from env vars HOME or USERPROFILE")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

func gitCmd(gitExe string, args []string) *exec.Cmd {
	if verbose {
		log.Println("Running:", gitExe, args)
	}
	cmd := exec.Command(gitExe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func findGit() (string, error) {
	return exec.LookPath("git")
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

func runOrExit(cmd *exec.Cmd) *exitStatusError {
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return &exitStatusError{exiterr, status.ExitStatus()}
			}
		} else {
			return &exitStatusError{exiterr, 1}
		}
	}

	return nil
}

func hasScheme(repoAddress string) bool {
	return strings.Contains(repoAddress, "://")
}

func parseRepo(repo, host string) (string, []string) {
	repo = strings.TrimSuffix(repo, ".git")

	var repoParts []string
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
		repo = repo + ".git"
		/// elif repo.count("/") == 1:
	} else {
		///	# Something from default host
		repoParts = []string{host}
		repoParts = append(repoParts, strings.Split(repo, "/")...)

		repo = fmt.Sprintf("git@%s:%s.git", host, repo)
		if verbose {
			log.Println("Did build repo:", repo)
		}
	}

	return repo, repoParts
}

/// Parse and split repo to path segments for repo address with scheme
/// https://github.com/some/thing.git -> []string{"github.com", "some", "thing"}
///
/// Input is guaranteed to contains "://"
func parseRepoWithScheme(repo string) []string {
	repo = strings.TrimSuffix(repo, ".git")
	parts := strings.SplitN(repo, "://", 2)
	return strings.Split(parts[1], "/")
}

func main() {
	if verbose {
		log.Println("Starting...")
	}
	gitExe, err := findGit()

	if err != nil {
		log.Fatalf("Could not find git executable: %s\n", err)
	}

	// if len(sys.argv) <= 1:
	// 	os.execlp("git", "git")

	if len(os.Args) <= 1 {
		exitErr := runOrExit(gitCmd(gitExe, []string{}))
		if nil != exitErr {
			exitErr.exit()
		}
	}

	/// repo = str(sys.argv[-1])
	repo := os.Args[len(os.Args)-1]

	/// if "/" not in repo:
	/// 	os.execlp("git", "git")
	if !strings.Contains(repo, "/") {
		exitErr := runOrExit(gitCmd(gitExe, []string{}))
		if exitErr != nil {
			exitErr.exit()
		}
	}

	host := os.Getenv("GIT_GET_HOST")
	if host == "" {
		host = "github.com"
	}

	verbose = (strings.TrimSpace(os.Getenv("GIT_GET_VERBOSE")) != "") || verbose

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
	if hasScheme(repo) {
		// https://github.com/some/thing -> []string{"github.com", "some", "thing"}
		repoParts = parseRepoWithScheme(repo)
	} else {
		repo, repoParts = parseRepo(repo, host)
	}

	targetDir := filepath.Join(repoParts...)
	targetDir = filepath.Join(path, targetDir)

	if verbose {
		log.Println("Using target dir:", targetDir)
	}
	err = os.MkdirAll(targetDir, os.ModePerm)
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
	exitErr := runOrExit(gitCmd(gitExe, []string{"clone", repo, targetDir}))

	if exitErr != nil {
		// Cleanup dir if empty

		// repo dir
		dir := targetDir
		rmDir(dir)

		// user dir
		dir = filepath.Dir(targetDir)
		rmDir(dir)

		// host dir
		dir = filepath.Dir(targetDir)
		rmDir(dir)
		exitErr.exit()
	}

	if verbose {
		log.Println("Ended...")
	}
}
