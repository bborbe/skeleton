// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import "testing"

// TestDeliberateFailure exists solely to produce a red default-branch CI run so
// the github-build-watcher webhook path can be verified end-to-end on real
// GitHub traffic. Removed immediately after the observation.
func TestDeliberateFailure(t *testing.T) {
	t.Fatal("deliberate failure: verifying github-build-watcher workflow_run webhook")
}
