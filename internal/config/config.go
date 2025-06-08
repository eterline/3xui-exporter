package config

import (
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
)

type Configuration struct {
	LogDir    string `arg:"--log-dir,env:LOG_DIR" help:"log file directory"`
	LogPretty bool   `arg:"--log-json,env:LOG_PRETTY" help:"log format in JSON syntax"`
	Debug     bool   `arg:"--debug,env:ENV_DEBUG" help:"allow debug logging level"`

	Listen     string `arg:"--listen" help:"Server listen address"`
	CrtFileSSL string `arg:"--certfile,env:CERT" help:"Server SSL certificate file"`
	KeyFileSSL string `arg:"--keyfile,env:KEY" help:"Server SSL key file"`

	DashboardURL      string `arg:"--url,env:URL" help:"3X-UI dashboard url"`
	DashboardBase     string `arg:"--base,env:BASE" help:"3X-UI dashboard url additional base"`
	DashboardLogin    string `arg:"--login,env:LOGIN" help:"3X-UI user login"`
	DashboardPassword string `arg:"--password,env:PASSWORD" help:"3X-UI user password"`
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
