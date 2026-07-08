package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// StreamSpammerEvents godoc
// @Id streamSpammerEvents
// @Summary Stream spammer lifecycle events
// @Tags Spammer
// @Description Streams spammer lifecycle events (create/update/status/membership/reorder/
// @Description delete) via Server-Sent Events so dashboards can update live. The stream is
// @Description unauthenticated and carries only safe metadata — no configs, seeds or logs.
// @Produce text/event-stream
// @Success 200 {string} string "SSE stream of spammer events"
// @Failure 500 {string} string "Streaming unsupported"
// @Router /api/spammers/events [get]
func (ah *APIHandler) StreamSpammerEvents(w http.ResponseWriter, r *http.Request) {
	// Intentionally no auth check: this stream only exposes safe metadata.

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	subscriptionID, eventChan := ah.daemon.SubscribeEvents()
	defer ah.daemon.UnsubscribeEvents(subscriptionID)

	// Heartbeat keeps proxies from closing the idle connection.
	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	// Prompt the client to do an initial sync as soon as the stream is established.
	fmt.Fprint(w, "event: ready\ndata: {}\n\n")
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		case <-heartbeat.C:
			fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		}
	}
}
