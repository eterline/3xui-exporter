package scrape

import (
	"context"
	"fmt"

	x3uiapi "github.com/eterline/x3ui-exporter/pkg/x3-ui-api"
)

type ClientStat struct {
	Name     string
	Protocol string
	Email    string
	Down     uint64
	Up       uint64
}

func (ctf ClientStat) DownTraffic() float64   { return float64(ctf.Down) }
func (ctf ClientStat) UpTraffic() float64     { return float64(ctf.Up) }
func (ctf ClientStat) EmailString() string    { return ctf.Email }
func (itf ClientStat) ProtocolString() string { return itf.Protocol }
func (itf ClientStat) NameString() string     { return itf.Name }

type ScraperXUI struct {
	api *x3uiapi.XUIClient
	ctx context.Context
}

func NewScraperXUI(ctx context.Context, c *x3uiapi.XUIClient) *ScraperXUI {
	return &ScraperXUI{
		api: c,
		ctx: ctx,
	}
}

func (scr *ScraperXUI) ScrapeInboundStats() ([]ClientStat, error) {
	inbs, err := scr.api.Inbounds(scr.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed fetch inbounds: %v", err)
	}

	stats := []ClientStat{}

	for _, inb := range inbs {
		if !statsAvailable(inb) {
			continue
		}

		name := inb.Remark
		proto := inb.Protocol

		for _, stat := range inb.ClientsStats {

			data := ClientStat{
				Name:     name,
				Protocol: proto,
				Email:    stat.Email,
				Up:       uint64(stat.Up),
				Down:     uint64(stat.Down),
			}

			stats = append(stats, data)
		}
	}

	return stats, nil
}

func statsAvailable(c x3uiapi.Inbound) bool {
	return len(c.ClientsStats) > 0
}
