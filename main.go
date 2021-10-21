package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/opensourceways/robot-gitee-plugin-lib/logrusutil"
	"github.com/opensourceways/robot-gitee-plugin-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-file-server/caller"
	"github.com/opensourceways/sync-file-server/grpc/server"
)

type options struct {
	endpoint          string
	platformTokenPath string
	platform          string
	port              string
	concurrentSize    int
}

func (o *options) addFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.port, "port", "8888", "Port to listen on.")
	fs.IntVar(&o.concurrentSize, "concurrent-size", 50, "the grpc server goroutine pool size.")
	fs.StringVar(&o.endpoint, "endpoint", "", "Path to server config file.(required)")
	fs.StringVar(&o.platform, "platform", "gitee", "Which code platform to provide rpc service and currently only supported gitte")
	fs.StringVar(&o.platformTokenPath, "platform-token-path", "/etc/platform/oauth", "Path to the file containing the Gitee OAuth secret.")
}

func (o options) validate() error {
	_, err := url.Parse(o.endpoint)
	return err
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

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.platformTokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}
	defer secretAgent.Stop()

	genToken := secretAgent.GetTokenGenerator(o.platformTokenPath)

	log := logrus.WithField("platform", o.platform)

	c := caller.New(o.platform, o.endpoint, genToken)
	if c == nil {
		log.Fatal(fmt.Printf("can't create caller instance with %s platform", o.platform))
	}

	if err := server.Start(":"+o.port, o.concurrentSize, c, log); err != nil {
		logrus.WithError(err).Fatal("error start grpc server.")
	}
}
