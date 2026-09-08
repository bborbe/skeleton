// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/go-skeleton/pkg/stats"
)

var _ = Describe("stats", func() {
	Describe("Mean", func() {
		It("computes the mean of a non-empty slice", func() {
			Expect(stats.Mean([]float64{1, 2, 3})).To(Equal(2.0))
		})

		It("returns zero for an empty slice", func() {
			Expect(stats.Mean(nil)).To(Equal(0.0))
		})
	})

	Describe("Median", func() {
		It("handles an odd-length slice", func() {
			Expect(stats.Median([]float64{3, 1, 2})).To(Equal(2.0))
		})

		It("averages the middle pair for even length", func() {
			Expect(stats.Median([]float64{1, 2, 3, 4})).To(Equal(2.5))
		})

		It("does not mutate the input", func() {
			input := []float64{3, 1, 2}
			_ = stats.Median(input)
			Expect(input).To(Equal([]float64{3, 1, 2}))
		})
	})

	Describe("Mode", func() {
		It("returns the most frequent value", func() {
			Expect(stats.Mode([]float64{1, 2, 2, 3})).To(Equal(2.0))
		})

		It("resolves ties to the smallest value", func() {
			Expect(stats.Mode([]float64{2, 1})).To(Equal(1.0))
		})
	})

	Describe("Variance and StdDev", func() {
		It("computes sample variance", func() {
			Expect(stats.Variance([]float64{2, 4, 4, 4, 5, 5, 7, 9})).To(BeNumerically("~", 4.571, 0.001))
		})

		It("returns zero for a single sample", func() {
			Expect(stats.Variance([]float64{5})).To(Equal(0.0))
		})

		It("std dev is the square root of variance", func() {
			Expect(stats.StdDev([]float64{2, 4, 4, 4, 5, 5, 7, 9})).To(BeNumerically("~", 2.138, 0.001))
		})
	})

	Describe("Range, Min, Max", func() {
		It("computes the range", func() {
			Expect(stats.Range([]float64{3, 1, 9})).To(Equal(8.0))
		})

		It("finds min and max", func() {
			Expect(stats.Min([]float64{3, 1, 9})).To(Equal(1.0))
			Expect(stats.Max([]float64{3, 1, 9})).To(Equal(9.0))
		})
	})

	Describe("Sum and Product", func() {
		It("sums the samples", func() {
			Expect(stats.Sum([]float64{1, 2, 3})).To(Equal(6.0))
		})

		It("products the samples", func() {
			Expect(stats.Product([]float64{2, 3, 4})).To(Equal(24.0))
		})

		It("empty product is one", func() {
			Expect(stats.Product(nil)).To(Equal(1.0))
		})
	})

	Describe("Quantile and Percentile", func() {
		It("computes the median as the 0.5 quantile", func() {
			Expect(stats.Quantile([]float64{1, 2, 3, 4}, 0.5)).To(Equal(2.5))
		})

		It("computes percentiles", func() {
			Expect(stats.Percentile([]float64{1, 2, 3, 4}, 25)).To(Equal(1.75))
		})

		It("clamps p to the endpoints", func() {
			Expect(stats.Quantile([]float64{1, 2, 3}, 0)).To(Equal(1.0))
			Expect(stats.Quantile([]float64{1, 2, 3}, 1)).To(Equal(3.0))
		})
	})

	Describe("MAD", func() {
		It("computes median absolute deviation", func() {
			Expect(stats.MAD([]float64{1, 2, 3, 4, 5})).To(Equal(1.0))
		})
	})

	Describe("Skewness and Kurtosis", func() {
		It("skewness of symmetric data is zero", func() {
			Expect(stats.Skewness([]float64{1, 2, 3, 4, 5, 6, 7})).To(BeNumerically("~", 0, 0.001))
		})

		It("kurtosis of small data is bounded", func() {
			Expect(math.Abs(stats.Kurtosis([]float64{1, 2, 3, 4, 5}))).To(BeNumerically("<", 10))
		})
	})

	Describe("Correlation", func() {
		It("perfectly correlated series give 1", func() {
			Expect(stats.Correlation([]float64{1, 2, 3}, []float64{2, 4, 6})).To(BeNumerically("~", 1, 0.001))
		})

		It("constant series give 0", func() {
			Expect(stats.Correlation([]float64{1, 1, 1}, []float64{2, 4, 6})).To(Equal(0.0))
		})
	})

	Describe("ZScore", func() {
		It("center value scores zero", func() {
			Expect(stats.ZScore(2, []float64{1, 2, 3})).To(BeNumerically("~", 0, 0.001))
		})
	})

	Describe("Trimmed and Winsorized mean", func() {
		It("trimmed mean drops the extremes", func() {
			Expect(stats.TrimmedMean([]float64{1, 2, 3, 4, 100}, 0.2)).To(BeNumerically("~", 3, 0.001))
		})

		It("winsorized mean clips outliers", func() {
			Expect(stats.WinsorizedMean([]float64{1, 2, 3, 4, 100}, 0.1)).To(BeNumerically("~", 4.8, 0.001))
		})
	})

	Describe("Confidence interval", func() {
		It("brackets the mean", func() {
			lo, hi := stats.ConfidenceInterval([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
			Expect(lo).To(BeNumerically("<", hi))
			Expect(lo).To(BeNumerically("<", 5.5))
			Expect(hi).To(BeNumerically(">", 5.5))
		})
	})

	Describe("Series transforms", func() {
		It("normalize sums to one", func() {
			Expect(stats.Sum(stats.Normalize([]float64{1, 3, 6}))).To(BeNumerically("~", 1, 0.001))
		})

		It("standardize has unit variance", func() {
			Expect(stats.Variance(stats.Standardize([]float64{1, 2, 3, 4, 5}))).To(BeNumerically("~", 1, 0.001))
		})

		It("cumulative sum is monotone", func() {
			Expect(stats.CumulativeSum([]float64{1, 2, 3})).To(Equal([]float64{1, 3, 6}))
		})

		It("differences invert the cumulative sum", func() {
			Expect(stats.Differences([]float64{1, 3, 6})).To(Equal([]float64{2, 3}))
		})

		It("lag shifts the series", func() {
			Expect(stats.Lag([]float64{1, 2, 3}, 1)).To(Equal([]float64{0, 1, 2}))
		})
	})
})
