package ksubdomain

import (
	"github.com/hktalent/ksubdomain/core/conf"
	"github.com/hktalent/ksubdomain/core/gologger"
	"github.com/urfave/cli/v2"
	"os"
)

func Main() {
	//os.Args = append([]string{""}, strings.Split("enum -d 5M -o superbet.ro.json -d superbet.ro -f /Users/51pwn/MyWork/scan4all/config/database/subdomain.txt", " ")...)
	app := &cli.App{
		Name:    conf.AppName,
		Version: conf.Version,
		Usage:   conf.Description,
		Commands: []*cli.Command{
			enumCommand,
			verifyCommand,
			testCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
}
