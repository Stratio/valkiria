package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/stratio/valkiria/routes"
	"golang.org/x/net/context"
	"os"
)

var (
	Version = "0.0.1"
	Build   = "HEAD"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Starting Valkiria...")
	log.Infof("Version: %v - Build: %v", Version, Build)
	serveCommand := cli.Command{
		Name:      "serve",
		ShortName: "s",
		Usage:     "Serve the API",
		Flags:     []cli.Flag{FlAddr},
		Action:    action(serveAction),
	}
	Run("valkiria", serveCommand)
}

func action(f func(c *cli.Context) error) func(c *cli.Context) {
	return func(c *cli.Context) {
		err := f(c)
		if err != nil {
			log.Error("PANIC: " + err.Error())
			os.Exit(1)
		}
	}
}

func serveAction(c *cli.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "addr", c.String("addr"))

	return ServeCmd(c, ctx, routes.Routes)
}
