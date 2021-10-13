package server

import (
	"context"

	"github.com/panjf2000/ants/v2"

	"github.com/opensourceways/sync-file-server/models"
	"github.com/opensourceways/sync-file-server/protocol"
)

func newSyncFileServer(concurrentSize int) (protocol.SyncFileServer, error) {
	p, err := ants.NewPool(concurrentSize, ants.WithOptions(ants.Options{
		PreAlloc:    true,
		Nonblocking: true,
	}))
	if err != nil {
		return nil, err
	}

	return &syncFileServer{pool: p}, nil
}

type syncFileServer struct {
	pool *ants.Pool
	protocol.UnimplementedSyncFileServer
}

func (s *syncFileServer) SyncFile(ctx context.Context, input *protocol.SyncFileRequest) (*protocol.Result, error) {
	b := input.Branch
	opt := models.SyncFileOption{
		Branch: models.Branch{
			Org:    b.Org,
			Repo:   b.Repo,
			Branch: b.Branch,
		},
		BranchSHA: b.BranchSha,
		Files:     input.Files,
	}

	err := s.pool.Submit(func() {
		if err := opt.Create(); err != nil {
			//log
		}
	})

	return new(protocol.Result), err
}

func (s *syncFileServer) SyncRepoFile(ctx context.Context, input *protocol.SyncRepoFileRequest) (*protocol.Result, error) {
	b := input.Branch
	opt := models.SyncRepoFileOption{
		Branch: models.Branch{
			Org:    b.Org,
			Repo:   b.Repo,
			Branch: b.Branch,
		},
		BranchSHA: b.BranchSha,
		FileNames: input.FileNames,
	}

	err := s.pool.Submit(func() {
		if err := opt.Create(); err != nil {
			//log
		}
	})

	return new(protocol.Result), err
}
