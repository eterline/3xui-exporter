package config

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	LogDir  string `arg:"--log-dir" help:"log file directory"`
	LogJSON bool   `arg:"--log-json" help:"log format in JSON syntax"`
	Debug   bool   `arg:"--debug" help:"allow debug logging level"`

	Listen     string `arg:"--listen" help:"Server listen address"`
	CrtFileSSL string `arg:"--certfile" help:"Server SSL certificate file"`
	KeyFileSSL string `arg:"--keyfile" help:"Server SSL key file"`

	DashboardURL      string `arg:"--url" help:"3X-UI dashboard url"`
	DashboardBase     string `arg:"--base" help:"3X-UI dashboard url additional base"`
	DashboardLogin    string `arg:"--login" help:"3X-UI user login"`
	DashboardPassword string `arg:"--password" help:"3X-UI user password"`
}

var (
	parserConfig = arg.Config{
		Program:           selfExec(),
		IgnoreEnv:         false,
		IgnoreDefault:     false,
		StrictSubcommands: true,
	}
)

func ParseArgs(c *Configuration) error {
	p, err := arg.NewParser(parserConfig, c)
	if err != nil {
		return err
	}

	err = p.Parse(os.Args[1:])
	if err == arg.ErrHelp {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	return err
}

func selfExec() string {
	exePath, err := os.Executable()
	if err != nil {
		return "monita"
	}

	return filepath.Base(exePath)
}
