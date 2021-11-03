package gitee

import (
	sdk "gitee.com/openeuler/go-gitee/gitee"
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

	repoNames := make([]string, 0, len(repos))
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Path)
	}
	return repoNames, nil
}

func (gp *giteePlatform) ListBranchesOfRepo(org, repo string) ([]backend.BranchInfo, error) {
	branches, err := gp.cli.GetRepoAllBranch(org, repo)
	if err != nil || len(branches) == 0 {
		return nil, err
	}

	transform := func(branch *sdk.Branch) backend.BranchInfo {
		sha := ""
		if branch.Commit != nil {
			sha = branch.Commit.Sha
		}

		return backend.BranchInfo{
			Name: branch.Name,
			SHA:  sha,
		}
	}

	infos := make([]backend.BranchInfo, 0, len(branches))
	for i := range branches {
		infos = append(infos, transform(&branches[i]))
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
		files[i] = backend.RepoFile{
			Path: item.Path,
			SHA:  item.Sha,
		}
	}
	return files, nil
}

func (gp *giteePlatform) GetFileConent(b backend.Branch, path string) (string, string, error) {
	content, err := gp.cli.GetPathContent(b.Org, b.Repo, path, b.Branch)
	if err != nil {
		return "", "", nil
	}
	return content.Sha, content.Content, nil
}
