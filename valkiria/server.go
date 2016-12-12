package main

import (
	"fmt"
	r "github.com/Stratio/valkiria/routes"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"net/http"
	"os"
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
	ro := r.NewRouter(ctx, routes)
	return http.ListenAndServe(c.String("ip"), ro)
}
