// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	libsentry "github.com/bborbe/sentry"
	"github.com/getsentry/sentry-go"
)

func NewSentryAlertHandler(sentryClient libsentry.Client) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		result := sentryClient.CaptureException(
			errors.Errorf(ctx, "my error"),
			&sentry.EventHint{
				Data: map[string]string{
					"key": "value",
				},
				Context: ctx,
			},
			nil,
		)
		_, _ = libhttp.WriteAndGlog(resp, "send sentry alert '%s' completed", *result)
	})
}
