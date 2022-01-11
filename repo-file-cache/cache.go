package cache

import (
	"github.com/opensourceways/repo-file-cache/models"
	"github.com/opensourceways/repo-file-cache/sdk"

	"github.com/opensourceways/sync-file-server/backend"
)

func NewBackendStorage(platform, endpoint string) backend.Storage {
	return &repoFileCache{
		cli:      sdk.NewSDK(endpoint, 3),
		platform: platform,
	}
}

type repoFileCache struct {
	cli      *sdk.SDK
	platform string
}

func (fc *repoFileCache) SaveFiles(b backend.Branch, branchSHA string, files []backend.File) error {
	opts := models.FileUpdateOption{
		Branch: models.Branch{
			Platform: fc.platform,
			Org:      b.Org,
			Repo:     b.Repo,
			Branch:   b.Branch,
		},
	}
	opts.BranchSHA = branchSHA

	n := len(files)
	fs := make([]models.File, n)
	for i := 0; i < n; i++ {
		item := &files[i]
		fs[i] = models.File{
			Path:    models.FilePath(item.Path),
			SHA:     item.SHA,
			Content: item.Content,
		}
	}
	opts.Files = fs

	return fc.cli.SaveFiles(opts)
}

func (fc *repoFileCache) GetFileSummary(b backend.Branch, fileName string) ([]backend.RepoFile, error) {
	v, err := fc.cli.GetFiles(
		models.Branch{
			Platform: fc.platform,
			Org:      b.Org,
			Repo:     b.Repo,
			Branch:   b.Branch,
		},
		fileName,
		true,
	)
	if err != nil {
		return nil, err
	}

	r := make([]backend.RepoFile, len(v.Files))

	for i, item := range v.Files {
		r[i].Path = item.Path.FullPath()
		r[i].SHA = item.SHA
	}

	return r, nil
}
