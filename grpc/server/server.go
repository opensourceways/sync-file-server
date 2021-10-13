package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/opensourceways/sync-file-server/backend"
	"github.com/opensourceways/sync-file-server/models"
	"github.com/opensourceways/sync-file-server/protocol"
)

var server *grpc.Server

func Start(port string, concurrentSize int, cli backend.Client) error {
	backend.RegisterClient(cli)

	if err := models.NewPool(concurrentSize * 10); err != nil {
		return err
	}

	syncFileServer, err := newSyncFileServer(concurrentSize)
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server = grpc.NewServer()

	protocol.RegisterSyncFileServer(server, syncFileServer)
	protocol.RegisterRepoServer(server, repoServer{})

	return server.Serve(listen)
}

func Stop() {
	if server != nil {
		server.Stop()
	}
}
