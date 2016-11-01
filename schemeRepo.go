package main

import "strings"

var schemeRepoParser repoParser = repoParser{
	test:  hasScheme,
	parse: parseSchemeRepo,
}

func hasScheme(repoAddress string) bool {
	return strings.Contains(repoAddress, "://")
}

/// Parse and split repo to path segments for repo address with scheme
/// https://github.com/some/thing.git -> []string{"github.com", "some", "thing"}
///
/// Input is guaranteed to contains "://"
func parseSchemeRepo(repo, host string) (string, []string) {
	repo = strings.TrimSuffix(repo, ".git")
	parts := strings.SplitN(repo, "://", 2)
	return repo, strings.Split(parts[1], "/")
}
