package models

import (
	"fmt"
	"path/filepath"
	"strings"

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

func fetchFileError(f string, err error) fetchFileResult {
	r := fetchFileResult{err: err}
	r.Path = f
	return r
}

func syncFile(branch Branch, branchSHA string, files []string) error {
	if len(files) == 0 {
		return nil
	}

	c := backend.GetClient()
	n := len(files)
	ch := make(chan fetchFileResult, n)

	task := func(f string) func() {
		return func() {
			sha, content, err := c.GetFileConent(branch, f)
			if err != nil {
				ch <- fetchFileError(f, err)
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
		}
	}

	for _, f := range files {
		if err := pool.Submit(task(f)); err != nil {
			ch <- fetchFileError(f, err)
		}
	}

	log := logEntryForBranch(branch, branchSHA)

	i := 0
	result := make([]File, 0, n)
	for r := range ch {
		if r.err == nil {
			result = append(result, r.File)
		} else {
			log.WithField("file", r.Path).WithError(r.err).Error("sync file")
		}

		if i++; i == n {
			break
		}
	}
	close(ch)

	if n := len(result); n > 0 {
		if err := c.SaveFiles(branch, branchSHA, result); err != nil {
			fs := make([]string, n)
			for i := range result {
				fs[i] = result[i].Path
			}
			return fmt.Errorf(
				"error to save files: %s, err: %v",
				strings.Join(fs, "; "), err,
			)
		}
	}
	return nil
}

func parseFileName(s string) string {
	return filepath.Base(s)
}

func logEntryForBranch(b Branch, sha string) *logrus.Entry {
	return logrus.WithFields(
		logrus.Fields{
			"org":        b.Org,
			"repo":       b.Repo,
			"branch":     b.Branch,
			"branch sha": sha,
		},
	)
}
