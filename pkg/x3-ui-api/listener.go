package x3uiapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUpdaterClosed = errors.New("stats updater closed")
)

type StatsMessage struct {
	Updates TrafficUpdates
	Err     error
}

type StatsHandle struct {
	statChan chan StatsMessage
}

func NewStatsHandler() *StatsHandle {
	return &StatsHandle{
		statChan: make(chan StatsMessage),
	}
}

func (sl *StatsHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stat := TrafficUpdates{}
	err := json.NewDecoder(r.Body).Decode(&stat)

	msg := StatsMessage{
		Updates: stat,
		Err:     nil,
	}

	if err != nil {
		msg.Err = fmt.Errorf("failed to decode traffic stats from request: %w", err)
	}

	sl.statChan <- msg
	w.WriteHeader(http.StatusNoContent)
}

func (sl *StatsHandle) Updates(ctx context.Context) <-chan StatsMessage {
	updateChan := make(chan StatsMessage)

	go func() {
		defer close(updateChan)

		for {
			select {

			case <-ctx.Done():
				return

			case stats, ok := <-sl.statChan:
				if !ok {
					updateChan <- StatsMessage{
						Updates: TrafficUpdates{},
						Err:     ErrUpdaterClosed,
					}
					return
				}
				updateChan <- stats
			}
		}
	}()

	return updateChan
}

func (sl *StatsHandle) Close() error {
	close(sl.statChan)
	return nil
}
