// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build go1.18
// +build go1.18

// Package metricdatatest provides testing functionality for use with the
// metricdata package.
package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Datatypes are the concrete data-types the metricdata package provides.
type Datatypes interface {
	metricdata.DataPoint[float64] |
		metricdata.DataPoint[int64] |
		metricdata.Gauge[float64] |
		metricdata.Gauge[int64] |
		metricdata.Histogram |
		metricdata.HistogramDataPoint |
		metricdata.Metrics |
		metricdata.ResourceMetrics |
		metricdata.ScopeMetrics |
		metricdata.Sum[float64] |
		metricdata.Sum[int64]

	// Interface types are not allowed in union types, therefore the
	// Aggregation and Value type from metricdata are not included here.
}

type config struct {
	ignoreTimestamp bool
}

// Option allows for fine grain control over how AssertEqual operates.
type Option interface {
	apply(cfg config) config
}

type fnOption func(cfg config) config

func (fn fnOption) apply(cfg config) config {
	return fn(cfg)
}

// IgnoreTimestamp disables checking if timestamps are different.
func IgnoreTimestamp() Option {
	return fnOption(func(cfg config) config {
		cfg.ignoreTimestamp = true
		return cfg
	})
}

// AssertEqual asserts that the two concrete data-types from the metricdata
// package are equal.
func AssertEqual[T Datatypes](t *testing.T, expected, actual T, opts ...Option) bool {
	t.Helper()

	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	// Generic types cannot be type asserted. Use an interface instead.
	aIface := interface{}(actual)

	var r []string
	switch e := interface{}(expected).(type) {
	case metricdata.DataPoint[int64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[int64]), cfg)
	case metricdata.DataPoint[float64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[float64]), cfg)
	case metricdata.Gauge[int64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[int64]), cfg)
	case metricdata.Gauge[float64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[float64]), cfg)
	case metricdata.Histogram:
		r = equalHistograms(e, aIface.(metricdata.Histogram), cfg)
	case metricdata.HistogramDataPoint:
		r = equalHistogramDataPoints(e, aIface.(metricdata.HistogramDataPoint), cfg)
	case metricdata.Metrics:
		r = equalMetrics(e, aIface.(metricdata.Metrics), cfg)
	case metricdata.ResourceMetrics:
		r = equalResourceMetrics(e, aIface.(metricdata.ResourceMetrics), cfg)
	case metricdata.ScopeMetrics:
		r = equalScopeMetrics(e, aIface.(metricdata.ScopeMetrics), cfg)
	case metricdata.Sum[int64]:
		r = equalSums(e, aIface.(metricdata.Sum[int64]), cfg)
	case metricdata.Sum[float64]:
		r = equalSums(e, aIface.(metricdata.Sum[float64]), cfg)
	default:
		// We control all types passed to this, panic to signal developers
		// early they changed things in an incompatible way.
		panic(fmt.Sprintf("unknown types: %T", expected))
	}

	if len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation, opts ...Option) bool {
	t.Helper()

	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if r := equalAggregations(expected, actual, cfg); len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}
