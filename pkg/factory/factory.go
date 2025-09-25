// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"net/http"

	libsentry "github.com/bborbe/sentry"

	"github.com/bborbe/skeleton/pkg/handler"
)

func CreateTestLoglevelHandler() http.Handler {
	return handler.NewTestLoglevelHandler()
}

func CreateSentryAlertHandler(sentryClient libsentry.Client) http.Handler {
	return handler.NewSentryAlertHandler(sentryClient)
}
