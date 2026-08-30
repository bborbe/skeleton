// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	libhttp "github.com/bborbe/http"
)

var processStart = time.Now()

func NewUptimeHandler() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add(libhttp.ContentTypeHeaderName, libhttp.ApplicationJSONContentType)
		uptime := time.Since(processStart)
		body, err := json.Marshal(struct {
			UptimeSeconds float64 `json:"uptime_seconds"`
			StartedAt     string  `json:"started_at"`
		}{
			UptimeSeconds: uptime.Seconds(),
			StartedAt:     processStart.Format(time.RFC3339),
		})
		if err != nil {
			http.Error(
				resp,
				fmt.Errorf("marshal uptime response: %v", err).Error(),
				http.StatusInternalServerError,
			)
			return
		}
		_, _ = resp.Write(body)
	})
}
