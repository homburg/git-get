package main

import "strings"

var sshRepoParser repoParser = repoParser{
	test: func(repo string) bool {
		return strings.Contains(repo, ":")
	},
	parse: parseSshRepo,
}

func parseSshRepo(repo, host string) (string, []string) {

	var repoParts []string

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

	return repo, repoParts
}
