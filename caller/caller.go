package caller

import (
	"github.com/opensourceways/sync-file-server/backend"
)

type Caller interface {
	backend.CodePlatform
	backend.Storage
}

func New(platform string, cacheServerURL string, getToken func() []byte) Caller {
	cacheCall := NewCacheCaller(cacheServerURL, platform)

	switch platform {
	case "gitee":
		return NewGiteeCaller(getToken, cacheCall)
	default:
		return nil
	}
}
