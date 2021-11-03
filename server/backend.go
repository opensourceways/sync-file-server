package main

import (
	"fmt"

	"github.com/opensourceways/sync-file-server/backend"
	"github.com/opensourceways/sync-file-server/repo-file-cache"
	"github.com/opensourceways/sync-file-server/server/gitee"
)

func newBackend(fileCacheEndpoint, platform string, token func() []byte) (backend.Client, error) {
	var cli backend.CodePlatform

	switch platform {
	case "gitee":
		cli = gitee.NewBackend(token)
	default:
		return nil, fmt.Errorf("unknown platform:%s", platform)
	}

	return struct {
		backend.CodePlatform
		backend.Storage
	}{
		CodePlatform: cli,
		Storage:      cache.NewBackendStorage(platform, fileCacheEndpoint),
	}, nil
}
