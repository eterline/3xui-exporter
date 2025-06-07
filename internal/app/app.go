package app

import (
	"context"
	"net/http"
	"time"

	"github.com/eterline/x3ui-exporter/internal/config"
	"github.com/eterline/x3ui-exporter/internal/server"
	"github.com/eterline/x3ui-exporter/internal/service/metrics"
	"github.com/eterline/x3ui-exporter/pkg/logger"
	"github.com/eterline/x3ui-exporter/pkg/toolkit"
	x3uiapi "github.com/eterline/x3ui-exporter/pkg/x3-ui-api"
	"github.com/go-chi/chi/v5"
)

var (
	waitDuration = 5 * time.Second
	log          logger.LogWorker
)

func Execute(cfg config.Configuration, root *toolkit.AppStarter) {
	log = logger.ReturnEntry()

	log.Info("service started")
	defer log.Info("service stopped")

	registry := metrics.NewMetricsReg(log)
	stats := x3uiapi.NewStatsHandler()
	defer stats.Close()

	go processUpdate(root.Context, stats, registry)
	go startServer(root.Context, cfg, registry, stats)

	log.Infof("server listen in: %s", cfg.Listen)

	root.Wait()
	root.WaitThreads(waitDuration)
}

func processUpdate(ctx context.Context, stats *x3uiapi.StatsHandle, reg *metrics.MetricsReg) {
	for update := range stats.Updates(ctx) {
		u := update
		log.Debug("got new traffic stats")

		if u.Err != nil {
			log.Errorf("failed update traffic stats: %v", u.Err)
			continue
		}

		for _, cl := range u.Updates.Client {
			reg.UpdateClient(cl)
		}

		for _, cl := range u.Updates.Inbound {
			reg.UpdateInbound(cl)
		}

	}
}

func startServer(ctx context.Context, cfg config.Configuration, reg *metrics.MetricsReg, stats *x3uiapi.StatsHandle) {

	r := chi.NewMux()
	r.Get("/metric", reg.Metric().ServeHTTP)
	r.Post("/metric", stats.ServeHTTP)

	srv := server.NewMetricsServer(r, cfg.Listen)
	go func() {
		err := srv.Listen(cfg.CrtFileSSL, cfg.KeyFileSSL)
		switch {
		case err == http.ErrServerClosed:
			return
		case err == nil:
			return
		default:
			log.Fatalf("server closed with error: %v", err)
		}
	}()

	<-ctx.Done()
	srv.Stop()
}
