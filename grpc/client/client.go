package client

import (
	"google.golang.org/grpc"

	"github.com/opensourceways/sync-file-server/protocol"
)

func NewClient(endpoint string) (Client, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	repoCli := protocol.NewRepoClient(conn)
	syncFileCli := protocol.NewSyncFileClient(conn)

	return &client{
		RepoClient:     repoCli,
		SyncFileClient: syncFileCli,
		conn:           conn,
	}, nil
}

type Client interface {
	protocol.RepoClient
	protocol.SyncFileClient

	Disconnect() error
}

type client struct {
	protocol.RepoClient
	protocol.SyncFileClient

	conn *grpc.ClientConn
}

func (c *client) Disconnect() error {
	return c.conn.Close()
}
