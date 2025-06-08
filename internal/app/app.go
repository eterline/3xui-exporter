package app

import (
	"context"
	"net/http"
	"time"

	"github.com/eterline/x3ui-exporter/internal/config"
	"github.com/eterline/x3ui-exporter/internal/server"
	"github.com/eterline/x3ui-exporter/internal/service/metrics"
	"github.com/eterline/x3ui-exporter/internal/service/scrape"
	"github.com/eterline/x3ui-exporter/pkg/logger"
	"github.com/eterline/x3ui-exporter/pkg/toolkit"
	x3uiapi "github.com/eterline/x3ui-exporter/pkg/x3-ui-api"
	"github.com/go-chi/chi/v5"
)

var (
	scrapeDuration = 15 * time.Second
	waitDuration   = 5 * time.Second
	log            logger.LogWorker
)

func Execute(cfg config.Configuration, root *toolkit.AppStarter) {
	log = logger.ReturnEntry()

	log.Info("service started")
	defer log.Info("service stopped")

	api, err := x3uiapi.NewClient(cfg.DashboardURL, cfg.DashboardBase, cfg.DashboardLogin, cfg.DashboardPassword, "")
	if err != nil {
		log.Fatalf("failed to init 3x-ui api: %v", err)
	}

	registry := metrics.NewMetricsReg(log)
	stats := x3uiapi.NewStatsHandler()
	defer stats.Close()

	scr := scrape.NewScraperXUI(root.Context, api)

	go processUpdate(root.Context, stats, registry)
	go processScrape(root.Context, scr, registry)
	go startServer(root.Context, cfg, registry, stats)

	log.Infof("server listen in: %s", cfg.Listen)

	root.Wait()
	root.WaitThreads(waitDuration)
}

func processUpdate(ctx context.Context, stats *x3uiapi.StatsHandle, reg *metrics.MetricsReg) {
	for u := range stats.Updates(ctx) {

		log.Debug("got new traffic stats")

		if u.Err != nil {
			log.Errorf("failed update traffic stats: %v", u.Err)
			continue
		}

		for _, cl := range u.Updates.Client {
			reg.UpdateClient(cl, cl)
		}

		for _, inb := range u.Updates.Inbound {
			reg.UpdateInbound(inb)
		}
	}
}

func processScrape(ctx context.Context, c *scrape.ScraperXUI, reg *metrics.MetricsReg) {
	ticker := time.NewTicker(scrapeDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			go func(scr *scrape.ScraperXUI, re *metrics.MetricsReg) {
				log.Debug("scrape stats started")
				defer log.Debug("stats scraped")

				stats, err := scr.ScrapeInboundStats()
				if err != nil {
					log.Error(err)
					return
				}

				for _, stat := range stats {
					re.UpdateStats(stat, stat)
				}
			}(c, reg)
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
