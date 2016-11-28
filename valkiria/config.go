package main

import "github.com/codegangsta/cli"

var (
	flIssuerURL = cli.StringFlag{
		Name:   "flag",
		Usage:  "flag ",
		Value:  "flag",
		EnvVar: "FLAG",
	}

	FlAddr = cli.StringFlag{
		Name:  "addr",
		Usage: "<ip>:<port> to listen on",
		Value: "127.0.0.1:8101",
	}
)