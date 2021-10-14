package models

import (
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-file-server/backend"
)

type SyncFileOption struct {
	Branch
	BranchSHA string   `json:"branch_sha" required:"true"`
	Files     []string `json:"files" required:"true"`
}

func (s SyncFileOption) Create() error {
	files := map[string][]string{}
	for _, f := range s.Files {
		name := parseFileName(f)
		files[name] = append(files[name], f)
	}

	for _, item := range files {
		syncFile(s.Branch, s.BranchSHA, item)
	}
	return nil
}

type fetchFileResult struct {
	File
	err error
}

func syncFile(branch Branch, branchSHA string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	c := backend.GetClient()
	n := len(files)
	ch := make(chan fetchFileResult, n)

	for _, f := range files {
		pool.Submit(func() {
			sha, content, err := c.GetFileConent(branch, f)
			if err != nil {
				ch <- fetchFileResult{
					File: File{
						RepoFile: RepoFile{
							Path: f,
						},
					},
					err: err,
				}
			} else {
				ch <- fetchFileResult{
					File: File{
						RepoFile: RepoFile{
							Path: f,
							SHA:  sha,
						},
						Content: content,
					},
				}
			}
		})
	}

	log := logrus.WithFields(
		logrus.Fields{
			"org":        branch.Org,
			"repo":       branch.Repo,
			"branch":     branch.Branch,
			"branch sha": branchSHA,
		},
	)
	i := 0
	result := make([]File, 0, n)
	for r := range ch {
		if r.err == nil {
			result = append(result, r.File)
		} else {
			log.WithField("file", r.File).WithError(r.err).Error("sync file")
		}

		if i++; i == n {
			break
		}
	}
	close(ch)

	if len(result) > 0 {
		return c.SaveFiles(branch, branchSHA, result)
	}
	return nil
}

func parseFileName(s string) string {
	return filepath.Base(s)
}
