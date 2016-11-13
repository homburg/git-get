package main

import (
	"fmt"
	"log"
	"strings"
)

var shorthandRepoParser repoParser = repoParser{
	test: func(repo string) bool {
		return !strings.Contains(repo, ":")
	},
	parse: parseShorthandRepo,
}

func parseShorthandRepo(repo, host string) (string, []string) {
	repoParts := []string{host}
	repoParts = append(repoParts, strings.Split(repo, "/")...)

	repo = fmt.Sprintf("git@%s:%s.git", host, repo)
	if verbose {
		log.Println("Did build repo:", repo)
	}

	return repo, repoParts
}
