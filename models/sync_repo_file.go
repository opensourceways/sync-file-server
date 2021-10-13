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

	for fileName, item := range files {
		todo, err := s.filterFile(fileName, item)
		if err != nil {

		}

		syncFile(s.Branch, s.BranchSHA, todo)
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

	i := 0
	n := len(files)
	todo := make([]string, 0, n)
	for _, item := range files {
		if item.SHA != m[item.Path] {
			todo[i] = item.Path
			i++
		}
	}
	return todo, nil
}
