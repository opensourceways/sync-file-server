package server

import (
	"context"

	"github.com/opensourceways/sync-file-server/backend"
	"github.com/opensourceways/sync-file-server/protocol"
)

type repoServer struct {
	protocol.UnimplementedRepoServer
}

func (rs repoServer) ListRepos(ctx context.Context, input *protocol.ListRepoRequest) (*protocol.ListRepoResponse, error) {
	c := backend.GetClient()

	v, err := c.ListRepos(input.GetOrg())
	if err != nil {
		return nil, err
	}

	return &protocol.ListRepoResponse{
		Repos: v,
	}, nil
}

func (rs repoServer) ListBranchesOfRepo(ctx context.Context, input *protocol.ListBranchesOfRepoRequest) (*protocol.ListBranchesOfRepoResponse, error) {
	c := backend.GetClient()

	v, err := c.ListBranchesOfRepo(input.GetOrg(), input.GetRepo())
	if err != nil {
		return nil, err
	}

	r := make([]*protocol.BranchInfo, 0, len(v))
	for _, item := range v {
		r = append(r, &protocol.BranchInfo{
			Name: item.Name,
			Sha:  item.SHA,
		})
	}
	return &protocol.ListBranchesOfRepoResponse{
		Branches: r,
	}, nil
}
