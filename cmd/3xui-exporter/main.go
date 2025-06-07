package main

import (
	"github.com/eterline/x3ui-exporter/internal/app"
	"github.com/eterline/x3ui-exporter/internal/config"
	"github.com/eterline/x3ui-exporter/pkg/logger"
	"github.com/eterline/x3ui-exporter/pkg/toolkit"
)

var (
	cfg = config.Configuration{
		LogDir:  "./logs",
		LogJSON: false,
		Debug:   false,

		Listen:     ":4500",
		CrtFileSSL: "",
		KeyFileSSL: "",
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
			logger.WithPrettyValue(!cfg.LogJSON),
			logger.WithPath(cfg.LogDir),
		)

		return err
	})

	app.Execute(cfg, root)
}
