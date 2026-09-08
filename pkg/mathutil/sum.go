// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathutil

// Sum returns the total of the given values.
// An empty slice sums to zero.
func Sum(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}
