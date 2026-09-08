// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathutil_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/go-skeleton/pkg/mathutil"
)

var _ = Describe("Sum", func() {
	It("sums a non-empty slice", func() {
		Expect(mathutil.Sum([]int{1, 2, 3})).To(Equal(6))
	})

	It("sums an empty slice to zero", func() {
		Expect(mathutil.Sum(nil)).To(Equal(0))
	})

	It("handles negatives", func() {
		Expect(mathutil.Sum([]int{-1, 1})).To(Equal(0))
	})
})
