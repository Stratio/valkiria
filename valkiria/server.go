package main

import (
	"github.com/codegangsta/cli"
	"fmt"
	"os"
	"net/http"
	"golang.org/x/net/context"
	r "github.com/stratio/valkiria/routes"
)

func Run(name string, commands ...cli.Command) {
	app := cli.NewApp()
	app.Name = name
	versionCommand := cli.Command{
		Name:      "version",
		ShortName: "v",
		Usage:     "Show the version",
		Action: func(_ *cli.Context) {
			fmt.Printf("%s-%s\n", Version, Build)
		},
	}
	app.Commands = append(commands, versionCommand)
	app.Version = fmt.Sprintf("%s-%s", Version, Build)
	app.Run(os.Args)
}

func ServeCmd(c *cli.Context, ctx context.Context, routes map[string]map[string]r.Handler) error {
	r := r.NewRouter(ctx, routes)
	return http.ListenAndServe(c.String("addr"), r)
}
