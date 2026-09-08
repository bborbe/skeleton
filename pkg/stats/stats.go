// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stats provides dependency-free descriptive statistics for
// float64 samples. All functions treat an empty input as undefined and
// return 0, matching the conventions of the mathutil package.
package stats

import (
	"math"
	"sort"
)

// Mean returns the arithmetic mean of the samples.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Median returns the middle value of a sorted copy of the samples.
func Median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	mid := n / 2
	if n%2 == 1 {
		return sorted[mid]
	}
	return (sorted[mid-1] + sorted[mid]) / 2
}

// Mode returns the most frequent value. Ties resolve to the smallest value.
func Mode(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	freq := make(map[float64]int)
	for _, v := range values {
		freq[v]++
	}
	best := values[0]
	bestCount := 0
	for v, c := range freq {
		if c > bestCount || (c == bestCount && v < best) {
			best = v
			bestCount = c
		}
	}
	return best
}

// Variance returns the sample variance (n-1 denominator).
func Variance(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}
	mean := Mean(values)
	var sum float64
	for _, v := range values {
		d := v - mean
		sum += d * d
	}
	return sum / float64(n-1)
}

// StdDev returns the square root of the sample variance.
func StdDev(values []float64) float64 {
	return math.Sqrt(Variance(values))
}

// Range returns the difference between the largest and smallest sample.
func Range(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min, max := values[0], values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

// Min returns the smallest sample.
func Min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the largest sample.
func Max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Sum returns the total of the samples.
func Sum(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum
}

// Product returns the product of the samples; an empty input yields 1.
func Product(values []float64) float64 {
	product := 1.0
	for _, v := range values {
		product *= v
	}
	return product
}

// GeometricMean returns the nth root of the product of n samples.
func GeometricMean(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	return math.Pow(Product(values), 1.0/float64(n))
}

// HarmonicMean returns the reciprocal of the mean of reciprocals.
func HarmonicMean(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		if v == 0 {
			return 0
		}
		sum += 1 / v
	}
	return float64(n) / sum
}

// Quantile returns the p-quantile (0 <= p <= 1) using linear interpolation.
func Quantile(values []float64, p float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[n-1]
	}
	idx := p * float64(n-1)
	lo := int(math.Floor(idx))
	hi := int(math.Ceil(idx))
	if lo == hi {
		return sorted[lo]
	}
	fraction := idx - float64(lo)
	return sorted[lo]*(1-fraction) + sorted[hi]*fraction
}

// Percentile returns the p-th percentile using Quantile(p/100).
func Percentile(values []float64, p float64) float64 {
	return Quantile(values, p/100)
}

// InterquartileRange returns Q3 - Q1.
func InterquartileRange(values []float64) float64 {
	return Quantile(values, 0.75) - Quantile(values, 0.25)
}

// MAD returns the median absolute deviation around the median.
func MAD(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	median := Median(values)
	deviations := make([]float64, 0, len(values))
	for _, v := range values {
		deviations = append(deviations, math.Abs(v-median))
	}
	return Median(deviations)
}

// Skewness returns the sample skewness (Fisher-Pearson).
func Skewness(values []float64) float64 {
	n := len(values)
	if n < 3 {
		return 0
	}
	mean := Mean(values)
	std := StdDev(values)
	if std == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		d := (v - mean) / std
		sum += d * d * d
	}
	return sum * float64(n) / float64((n-1)*(n-2))
}

// Kurtosis returns the excess sample kurtosis.
func Kurtosis(values []float64) float64 {
	n := len(values)
	if n < 4 {
		return 0
	}
	mean := Mean(values)
	std := StdDev(values)
	if std == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		d := (v - mean) / std
		sum += d * d * d * d
	}
	return sum*float64(n*(n+1))/float64((n-1)*(n-2)*(n-3)) - 3*float64((n-1)*(n-1))/float64((n-2)*(n-3))
}

// Covariance returns the sample covariance of paired samples.
func Covariance(xs, ys []float64) float64 {
	n := len(xs)
	if n < 2 || len(ys) != n {
		return 0
	}
	mx, my := Mean(xs), Mean(ys)
	var sum float64
	for i := 0; i < n; i++ {
		sum += (xs[i] - mx) * (ys[i] - my)
	}
	return sum / float64(n-1)
}

// Correlation returns the Pearson correlation coefficient.
func Correlation(xs, ys []float64) float64 {
	sx, sy := StdDev(xs), StdDev(ys)
	if sx == 0 || sy == 0 {
		return 0
	}
	return Covariance(xs, ys) / (sx * sy)
}

// ZScore returns how many standard deviations v is from the mean.
func ZScore(v float64, values []float64) float64 {
	std := StdDev(values)
	if std == 0 {
		return 0
	}
	return (v - Mean(values)) / std
}

// TrimmedMean returns the mean after dropping the lowest and highest p fraction.
func TrimmedMean(values []float64, p float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if p < 0 || p >= 0.5 {
		return Mean(values)
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	drop := int(math.Floor(p * float64(n)))
	trimmed := sorted[drop : n-drop]
	if len(trimmed) == 0 {
		return Mean(values)
	}
	return Mean(trimmed)
}

// WinsorizedMean replaces outliers with the p-quantile before averaging.
func WinsorizedMean(values []float64, p float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	lo, hi := Quantile(values, p), Quantile(values, 1-p)
	clipped := make([]float64, 0, n)
	for _, v := range values {
		switch {
		case v < lo:
			clipped = append(clipped, lo)
		case v > hi:
			clipped = append(clipped, hi)
		default:
			clipped = append(clipped, v)
		}
	}
	return Mean(clipped)
}

// StandardError returns the standard error of the mean.
func StandardError(values []float64) float64 {
	n := len(values)
	if n < 2 {
		return 0
	}
	return StdDev(values) / math.Sqrt(float64(n))
}

// ConfidenceInterval returns (lo, hi) for a z-based 95% interval.
func ConfidenceInterval(values []float64) (float64, float64) {
	se := StandardError(values)
	mean := Mean(values)
	return mean - 1.96*se, mean + 1.96*se
}

// N returns the number of samples.
func N(values []float64) int {
	return len(values)
}

// IsEmpty reports whether the sample set has no elements.
func IsEmpty(values []float64) bool {
	return len(values) == 0
}

// Clone returns a copy of the samples.
func Clone(values []float64) []float64 {
	return append([]float64(nil), values...)
}

// Normalize scales samples so they sum to 1.
func Normalize(values []float64) []float64 {
	sum := Sum(values)
	if sum == 0 {
		return append([]float64(nil), values...)
	}
	out := make([]float64, 0, len(values))
	for _, v := range values {
		out = append(out, v/sum)
	}
	return out
}

// Center subtracts the mean from each sample.
func Center(values []float64) []float64 {
	mean := Mean(values)
	out := make([]float64, 0, len(values))
	for _, v := range values {
		out = append(out, v-mean)
	}
	return out
}

// Standardize centers and scales to unit variance.
func Standardize(values []float64) []float64 {
	std := StdDev(values)
	if std == 0 {
		return Center(values)
	}
	out := make([]float64, 0, len(values))
	for _, v := range values {
		out = append(out, (v-Mean(values))/std)
	}
	return out
}

// CumulativeSum returns the running total at each position.
func CumulativeSum(values []float64) []float64 {
	out := make([]float64, 0, len(values))
	var running float64
	for _, v := range values {
		running += v
		out = append(out, running)
	}
	return out
}

// Differences returns the successive first differences.
func Differences(values []float64) []float64 {
	if len(values) < 2 {
		return nil
	}
	out := make([]float64, 0, len(values)-1)
	for i := 1; i < len(values); i++ {
		out = append(out, values[i]-values[i-1])
	}
	return out
}

// Lag shifts the series by n positions, filling the head with zeroes.
func Lag(values []float64, n int) []float64 {
	if n <= 0 {
		return Clone(values)
	}
	out := make([]float64, len(values))
	copy(out[n:], values[:len(values)-n])
	return out
}
