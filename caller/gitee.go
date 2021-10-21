package caller

import (
	"gitee.com/openeuler/go-gitee/gitee"
	gclient "github.com/opensourceways/robot-gitee-plugin-lib/giteeclient"

	"github.com/opensourceways/sync-file-server/backend"
)

type giteeClient interface {
	GetRepos(org string) ([]gitee.Project, error)
	GetRepoAllBranch(owner, repo string) ([]gitee.Branch, error)
	GetPathContent(org, repo, path, ref string) (gitee.Content, error)
	GetDirectoryTree(org, repo, sha string, recursive int32) (gitee.Tree, error)
}

type giteeCaller struct {
	giteeCli giteeClient
	CacheCaller
}

func (gc *giteeCaller) ListRepos(org string) ([]string, error) {
	repos, err := gc.giteeCli.GetRepos(org)
	if err != nil || len(repos) == 0 {
		return nil, err
	}

	repoNames := make([]string, 0, len(repos))
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Path)
	}
	return repoNames, nil
}

func (gc *giteeCaller) ListBranchesOfRepo(org, repo string) ([]backend.BranchInfo, error) {
	branches, err := gc.giteeCli.GetRepoAllBranch(org, repo)
	if err != nil || len(branches) == 0 {
		return nil, err
	}

	transform := func(branch gitee.Branch) backend.BranchInfo {
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
		infos = append(infos, transform(branches[i]))
	}
	return infos, err
}

func (gc *giteeCaller) ListAllFilesOfRepo(b backend.Branch) ([]backend.RepoFile, error) {
	trees, err := gc.giteeCli.GetDirectoryTree(b.Org, b.Repo, b.Branch, 1)
	if err != nil || len(trees.Tree) == 0 {
		return nil, err
	}

	transform := func(tree gitee.TreeBasic) backend.RepoFile {
		return backend.RepoFile{
			Path: tree.Path,
			SHA:  tree.Sha,
		}
	}

	files := make([]backend.RepoFile, 0, len(trees.Tree))
	for i := range trees.Tree {
		files = append(files, transform(trees.Tree[i]))
	}
	return files, nil
}

func (gc *giteeCaller) GetFileConent(b backend.Branch, path string) (string, string, error) {
	content, err := gc.giteeCli.GetPathContent(b.Org, b.Repo, path, b.Branch)
	if err != nil {
		return "", "", nil
	}
	return content.Sha, content.Content, nil
}

// NewGiteeCaller create gitee platform caller instance
func NewGiteeCaller(getToken func() []byte, cacheCaller CacheCaller) Caller {
	gcli := gclient.NewClient(getToken)
	return &giteeCaller{
		giteeCli:    gcli,
		CacheCaller: cacheCaller,
	}
}
