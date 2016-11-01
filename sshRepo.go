package main

import "strings"

var sshRepoParser repoParser = repoParser{
	test: func(repo string) bool {
		return strings.Contains(repo, ":")
	},
	parse: parseSshRepo,
}

func parseSshRepo(repo, host string) (string, []string) {

	// <user>@github.com/...
	hasUser := strings.Contains(repo, "@")

	var repoWithoutUser string
	if hasUser {
		repoWithoutUser = repo[strings.Index(repo, "@")+1 : len(repo)]
	} else {
		repoWithoutUser = repo
		repo = "git@" + repo
	}

	repoPath := strings.Replace(repoWithoutUser, ":", "/", -1)

	repoParts := strings.Split(repoPath, "/")

	repo = repo + ".git"

	return repo, repoParts
}
