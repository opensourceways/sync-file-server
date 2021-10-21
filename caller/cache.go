package caller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/opensourceways/sync-file-server/backend"
	"github.com/opensourceways/sync-file-server/caller/requests"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json;charset=UTF-8"
)

type file struct {
	Path    string `json:"path"`
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

type filesInfo struct {
	BranchSHA string `json:"branch_sha"`
	Files     []file `json:"files"`
}

type uploadOption struct {
	filesInfo
	Platform string `json:"platform" required:"true"`
	Org      string `json:"org" required:"true"`
	Repo     string `json:"repo" required:"true"`
	Branch   string `json:"branch" required:"true"`
}

type filesResult struct {
	Data filesInfo `json:"data"`
}

// CacheCaller is the caller of the file cache service
type CacheCaller struct {
	endpoint string
	platform string
}

func (cc *CacheCaller) SaveFiles(b backend.Branch, branchSHA string, files []backend.File) error {
	endpoint, err := cc.getCacheFileServerURL()
	if err != nil {
		return err
	}
	param := genUploadParam(b, files, branchSHA, cc.platform)

	jBody, err := json.Marshal(&param)
	if err != nil {
		return err
	}

	rst := requests.New(endpoint.String()).
		WithContext(context.Background()).
		WithMethod(http.MethodPost).
		SetHeader(contentType, contentTypeJson).
		WithBody(bytes.NewBuffer(jBody)).Do()

	if rst.Error() != nil {
		return rst.Error()
	}
	if rst.StatusCode() == http.StatusCreated {
		return nil
	}
	return fmt.Errorf("save files error:%s", string(rst.Body()))
}

func (cc *CacheCaller) GetFileSummary(b backend.Branch, fileName string) ([]backend.RepoFile, error) {
	endpoint, err := cc.urlOfGetFileCacheServer(b, fileName)
	if err != nil {
		return nil, err
	}
	var files filesResult
	err = requests.New(endpoint.String()).
		WithContext(context.Background()).
		SetHeader(contentType, contentTypeJson).
		Do().
		UnmarshalInto(&files)
	if err != nil {
		return nil, err
	}

	return resultToRepoFiles(files.Data), nil
}

func (cc *CacheCaller) getCacheFileServerURL() (*url.URL, error) {
	return url.Parse(cc.endpoint)
}

func (cc *CacheCaller) urlOfGetFileCacheServer(b backend.Branch, fileName string) (*url.URL, error) {
	serverURL, err := cc.getCacheFileServerURL()
	if err != nil {
		return nil, err
	}
	getFileUrl := &url.URL{
		Scheme: serverURL.Scheme,
		Host:   serverURL.Host,
		Path:   path.Join(serverURL.Path, cc.platform, b.Org, b.Repo, b.Branch, fileName),
	}
	return getFileUrl, nil
}

func resultToRepoFiles(info filesInfo) []backend.RepoFile {
	fl := len(info.Files)
	if fl == 0 {
		return []backend.RepoFile{}
	}

	rFiles := make([]backend.RepoFile, 0, fl)
	for _, v := range info.Files {
		rFiles = append(rFiles, backend.RepoFile{
			Path: v.Path,
			SHA:  v.SHA,
		})
	}
	return rFiles
}

func genUploadParam(b backend.Branch, files []backend.File, sha, platform string) uploadOption {
	fs := make([]file, 0, len(files))
	for _, v := range files {
		fs = append(fs, file{
			Path:    v.Path,
			SHA:     v.SHA,
			Content: v.Content,
		})
	}
	fInfo := filesInfo{
		BranchSHA: sha,
		Files:     fs,
	}
	return uploadOption{
		filesInfo: fInfo,
		Platform:  platform,
		Org:       b.Org,
		Repo:      b.Repo,
		Branch:    b.Branch,
	}
}

// NewCacheCaller create a caller instance of the file cache service
func NewCacheCaller(endpoint, platform string) CacheCaller {
	return CacheCaller{endpoint: endpoint, platform: platform}
}
