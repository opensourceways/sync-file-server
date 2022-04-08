package gitee

import (
	"github.com/opensourceways/community-robot-lib/giteeclient"

	"github.com/opensourceways/sync-file-server/backend"
)

func NewPlatform(getToken func() []byte) backend.CodePlatform {
	return &giteePlatform{
		cli: giteeclient.NewClient(getToken),
	}
}

type giteePlatform struct {
	cli giteeclient.Client
}

func (gp *giteePlatform) ListRepos(org string) ([]string, error) {
	repos, err := gp.cli.GetRepos(org)
	if err != nil || len(repos) == 0 {
		return nil, err
	}

	repoNames := make([]string, len(repos))

	for i := range repos {
		repoNames[i] = repos[i].Path
	}

	return repoNames, nil
}

func (gp *giteePlatform) ListBranchesOfRepo(org, repo string) ([]backend.BranchInfo, error) {
	branches, err := gp.cli.GetRepoAllBranch(org, repo)
	if err != nil || len(branches) == 0 {
		return nil, err
	}

	infos := make([]backend.BranchInfo, len(branches))

	for i := range branches {
		item := &branches[i]

		infos[i].Name = item.GetName()
		infos[i].SHA = item.GetCommit().GetSha()
	}

	return infos, err
}

func (gp *giteePlatform) ListAllFilesOfRepo(b backend.Branch) ([]backend.RepoFile, error) {
	trees, err := gp.cli.GetDirectoryTree(b.Org, b.Repo, b.Branch, 1)
	if err != nil || len(trees.Tree) == 0 {
		return nil, err
	}

	files := make([]backend.RepoFile, len(trees.Tree))

	for i := range trees.Tree {
		item := &trees.Tree[i]

		files[i].Path = item.Path
		files[i].SHA = item.Sha
	}

	return files, nil
}

func (gp *giteePlatform) GetFileConent(b backend.Branch, path string) (string, string, error) {
	content, err := gp.cli.GetPathContent(b.Org, b.Repo, path, b.Branch)
	if err != nil {
		return "", "", err
	}
	return content.Sha, content.Content, nil
}
