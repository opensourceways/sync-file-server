package server

import (
	"net"

	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/opensourceways/sync-file-server/backend"
	"github.com/opensourceways/sync-file-server/models"
	"github.com/opensourceways/sync-file-server/protocol"
)

func Start(port string, concurrentSize int, cli backend.Client) error {
	clears := []func(){
		models.Stop,
	}
	defer func() {
		for _, f := range clears {
			f()
		}
		logrus.Info("server exits.")
	}()

	backend.RegisterClient(cli)

	if err := models.NewPool(concurrentSize*10, logWapper{}); err != nil {
		return err
	}

	syncFileServer, err := newSyncFileServer(concurrentSize)
	if err != nil {
		return err
	}
	clears = append(clears, syncFileServer.Stop)

	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	protocol.RegisterSyncFileServer(server, syncFileServer)
	protocol.RegisterRepoServer(server, repoServer{})

	run(server, listen)
	return nil
}

func run(server *grpc.Server, listen net.Listener) {
	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		logrus.Errorf("grpc server exit...")
		server.Stop()
	})

	if err := server.Serve(listen); err != nil {
		logrus.Error(err)
	}
}

type logWapper struct{}

func (l logWapper) Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)

}
