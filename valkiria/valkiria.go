package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/Stratio/valkiria/routes"
	"golang.org/x/net/context"
	"os"
)

var (
	Version string
	Build   string
)

func main() {
	serveCommand := cli.Command{
		Name:      "agent",
		ShortName: "a",
		Usage:     "Serve Agent API",
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
	level, err := log.ParseLevel(c.String("log"))
	if err != nil {
		log.Error("PANIC: Level not found. " + err.Error())
		os.Exit(1)
	}
	log.SetLevel(level)
	return ServeCmd(c, ctx, routes.Routes)
}
