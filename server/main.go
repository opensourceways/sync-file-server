package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	"github.com/opensourceways/community-robot-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-file-server/grpc/server"
)

type options struct {
	port              string
	endpoint          string
	platform          string
	platformTokenPath string
	concurrentSize    int
}

func (o *options) addFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.port, "port", "8888", "Port to listen on.")
	fs.StringVar(&o.endpoint, "endpoint", "", "The endpoint of repo file cache")
	fs.StringVar(&o.platform, "platform", "gitee", "The code platform which implements rpc service. Currently only gitee is supported")
	fs.StringVar(&o.platformTokenPath, "platform-token-path", "/etc/platform/oauth", "The path to the token file which is used to access code platform.")

	fs.IntVar(&o.concurrentSize, "concurrent-size", 2000, "The grpc server goroutine pool size.")
}

func (o *options) validate() error {
	v, err := url.Parse(o.endpoint)
	if err != nil {
		return err
	}
	o.endpoint = v.String()

	if o.concurrentSize <= 0 {
		return fmt.Errorf("concurrent size must be bigger than 0")
	}
	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options
	o.addFlags(fs)
	_ = fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("sync-file-server")

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	logrus.WithField("platform", o.platform).Info("Starts sync file server.")

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.platformTokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error to start secret agent.")
	}
	defer secretAgent.Stop()

	getToken := secretAgent.GetTokenGenerator(o.platformTokenPath)

	backend, err := newBackend(o.endpoint, o.platform, getToken)
	if err != nil {
		logrus.WithError(err).Fatal("Error to generate backend")
	}

	if err := server.Start(":"+o.port, o.concurrentSize, backend); err != nil {
		logrus.WithError(err).Fatal("Error to start grpc server.")
	}
}
