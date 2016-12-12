package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Stratio/valkiria/dbus"
	"github.com/Stratio/valkiria/routes"
	"github.com/codegangsta/cli"
	"golang.org/x/net/context"
	"os"
)

var (
	Version string
	Build   string
)

func main() {
	agentCommand := cli.Command{
		Name:      "agent",
		ShortName: "a",
		Usage:     "Serve Agent API",
		Flags:     []cli.Flag{FlAddr, FlLog},
		Action:    action(agentAction),
	}
	orchestratorCommand := cli.Command{
		Name:      "orchestrator",
		ShortName: "o",
		Usage:     "Serve orchestrator API",
		Flags:     []cli.Flag{FlAddr, FlLog},
		Action:    action(orchestratorAction),
	}
	Run("valkiria", orchestratorCommand, agentCommand)
}

func action(f func(c *cli.Context) error) func(c *cli.Context) {
	return func(c *cli.Context) {
		setLogLevel(c)
		err := f(c)
		if err != nil {
			log.Error("PANIC: " + err.Error())
			os.Exit(1)
		}
	}
}

func agentAction(c *cli.Context) error {
	setDBusInstance()
	return startServer(c)
}

func orchestratorAction(c *cli.Context) error {
	return startServer(c)
}

func startServer(c *cli.Context) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "ip", c.String("ip"))
	log.Info("Serving api")
	return ServeCmd(c, ctx, routes.RoutesOrchestrator)
}
func setLogLevel(c *cli.Context) {
	level, err := log.ParseLevel(c.String("log"))
	if err != nil {
		log.Fatalf("Wrong log level: %v. Use INFO or DEBUG options", err.Error())
		os.Exit(1)
	}
	log.SetLevel(level)
	log.Infof("Log is ready in level: %v", level)
}

func setDBusInstance() {
	if err := dbus.DbusInstance.NewDBus(); err != nil {
		log.Fatalf("Error initializating D-Bus system. Stop the program. FATAL: %v", err)
		os.Exit(1)
	}
	log.Info("DBus is ready")
}
