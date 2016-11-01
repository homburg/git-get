package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var verbose = false

type repoParser struct {
	test  func(string) bool
	parse func(string, string) (string, []string)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type repoParserList []repoParser

/// Apply the first matching repo parser
func (r repoParserList) parse(repo, host string) (string, []string) {
	for _, parser := range r {
		if parser.test(repo) {
			return parser.parse(repo, host)
		}
	}

	return repo, []string{}
}

/// Prioritized list of repo parsers
var repoParsers = repoParserList([]repoParser{
	schemeRepoParser,
	shorthandRepoParser,
	sshRepoParser,
})

func main() {
	if verbose {
		log.Println("Starting...")
	}
	gitExe, err := findGit()

	if err != nil {
		log.Fatalf("Could not find git executable: %s\n", err)
	}

	// If no argument given, just run git and quit
	if len(os.Args) <= 1 {
		exitErr := runOrExit(gitCmd(gitExe, []string{}))
		if nil != exitErr {
			exitErr.exit()
		}
	}

	repo := os.Args[len(os.Args)-1]

	// Repo address must contains at least one "/"
	// otherwise run git and quit
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
		path = defaultPath()
	}

	if verbose {
		log.Println("Using base:", path)
	}

	repo = strings.TrimSuffix(repo, ".git")

	// Find the right parser for the repo address
	repo, repoParts := repoParsers.parse(repo, host)

	targetDir := filepath.Join(repoParts...)
	targetDir = filepath.Join(path, targetDir)

	if verbose {
		log.Println("Using target dir:", targetDir)
	}

	err = os.MkdirAll(targetDir, os.ModePerm)
	if nil != err {
		log.Fatal(err)
	}

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
