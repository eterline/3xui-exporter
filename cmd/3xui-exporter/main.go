package main

import (
	"github.com/eterline/x3ui-exporter/internal/app"
	"github.com/eterline/x3ui-exporter/internal/config"
	"github.com/eterline/x3ui-exporter/pkg/logger"
	"github.com/eterline/x3ui-exporter/pkg/toolkit"
)

var (
	cfg = config.Configuration{
		LogDir:    "./logs",
		LogPretty: true,
		Debug:     false,

		Listen:     ":4500",
		CrtFileSSL: "",
		KeyFileSSL: "",

		DashboardURL:      "",
		DashboardBase:     "",
		DashboardLogin:    "",
		DashboardPassword: "",
	}
)

func main() {
	root := toolkit.InitAppStart(func() error {
		var err error

		err = config.ParseArgs(&cfg)
		if err != nil {
			return err
		}

		err = logger.InitLogger(
			logger.WithDevEnvBool(cfg.Debug),
			logger.WithPrettyValue(cfg.LogPretty),
			logger.WithPath(cfg.LogDir),
		)

		return err
	})

	app.Execute(cfg, root)
}
