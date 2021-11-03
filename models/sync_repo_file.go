package models

import (
	"github.com/opensourceways/sync-file-server/backend"
)

type File = backend.File
type Branch = backend.Branch
type RepoFile = backend.RepoFile

type SyncRepoFileOption struct {
	Branch
	BranchSHA string   `json:"branch_sha" required:"true"`
	FileNames []string `json:"file_names" required:"true"`
}

func (s SyncRepoFileOption) Create() error {
	c := backend.GetClient()

	allFiles, err := c.ListAllFilesOfRepo(s.Branch)
	if err != nil {
		return err
	}

	files := make(map[string][]RepoFile)
	for _, f := range s.FileNames {
		files[f] = []RepoFile{}
	}

	for _, f := range allFiles {
		name := parseFileName(f.Path)
		if v, ok := files[name]; ok {
			files[name] = append(v, f)
		}
	}

	log := logEntryForBranch(s.Branch, s.BranchSHA)
	logError := func(f, msg string, err error) {
		log.WithField("file name", f).WithError(err).Error(msg)
	}

	for fileName, item := range files {
		todo, err := s.filterFile(fileName, item)
		if err != nil {
			logError(fileName, "filter file", err)
			return err
		}

		if err := syncFile(s.Branch, s.BranchSHA, todo); err != nil {
			logError(fileName, "sync file", err)
		}
	}
	return nil
}

func (s SyncRepoFileOption) filterFile(fileName string, files []RepoFile) ([]string, error) {
	if len(files) == 0 {
		return nil, nil
	}

	c := backend.GetClient()

	summary, err := c.GetFileSummary(s.Branch, fileName)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	for _, item := range summary {
		m[item.Path] = item.SHA
	}

	todo := make([]string, 0, len(files))
	for _, item := range files {
		if item.SHA != m[item.Path] {
			todo = append(todo, item.Path)
		}
	}
	return todo, nil
}
