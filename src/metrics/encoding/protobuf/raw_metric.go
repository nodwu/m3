// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package protobuf

import (
	"github.com/m3db/m3/src/metrics/generated/proto/metricpb"
	"github.com/m3db/m3/src/metrics/metric/aggregated"
	"github.com/m3db/m3/src/metrics/metric/id"
)

var r aggregated.RawMetric = newRawMetric()

type rawMetric struct {
	pb metricpb.ChunkedMetricWithStoragePolicyWithEncodeTime
}

func newRawMetric() *rawMetric {
	return &rawMetric{}
}

// ID is the metric identifier.
func (m rawMetric) ID() (id.RawID, error) {
	return nil, nil
}

// TimeNanos is the metric timestamp in nanoseconds.
func (m rawMetric) TimeNanos() (int64, error) {
	return 0, nil
}

// Value is the metric value.
func (m rawMetric) Value() (float64, error) {
	return 0, nil
}

// Metric is the metric object represented by the raw metric.
func (m rawMetric) Metric() (aggregated.Metric, error) {
	return aggregated.Metric{}, nil
}

// Bytes are the bytes backing this raw metric.
func (m rawMetric) Bytes() []byte {
	return nil
}

// Reset resets the raw data.
func (m *rawMetric) Reset(data []byte) {
}
