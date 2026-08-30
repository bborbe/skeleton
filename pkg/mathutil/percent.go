// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathutil

// Percent returns what percentage part is of total, rounded down.
func Percent(part, total int) int {
	return part * 100 / total
}
