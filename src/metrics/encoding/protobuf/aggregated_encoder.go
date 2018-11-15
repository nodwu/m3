// Copyright (c) 2018 Uber Technologies, Inc.
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
	"encoding/binary"
	"fmt"
	"math"

	"github.com/m3db/m3/src/metrics/encoding/msgpack"
	"github.com/m3db/m3/src/metrics/generated/proto/metricpb"
	"github.com/m3db/m3/src/metrics/metric/aggregated"
	"github.com/m3db/m3x/pool"
)

type Encoder interface {
	EncodeChunkedMetricAndEncodeTime(
		cmp aggregated.ChunkedMetricWithStoragePolicy,
		encodedAtNanos int64,
	) error

	Encode(m aggregated.MetricWithStoragePolicy, encodedAtNanos int64) error
	Buffer() msgpack.Buffer
}

// AggregatedEncoder encodes aggregated metrics.
type AggregatedEncoder interface {
	// Len returns the number of bytes accumulated in the encoder so far.
	Len() int

	// Reset resets the encoder buffer with optional initial data.
	Reset(initData []byte)

	// Truncate discards all but the first n encoded bytes but continues to use
	// the same allocated storage. If n is negative or greater than the length of
	// encoded buffer, an error is returned.
	Truncate(n int) error

	// EncodeChunkedMetricWithStoragePolicyAndEncodeTime encodes a chunked metric with
	// an applicable storage policy, alongside the time at which encoding happens.
	EncodeChunkedMetricWithStoragePolicyAndEncodeTime(
		cmp aggregated.ChunkedMetricWithStoragePolicy,
		encodedAtNanos int64,
	) error

	// Relinquish relinquishes ownership of the encoded byte stream to the caller,
	// and resets the internal encoding buffer.
	Relinquish() Buffer
}

type aggregatedEncoder struct {
	pool           pool.BytesPool
	initBufSize    int
	maxMessageSize int

	cm   metricpb.ChunkedMetricWithStoragePolicyWithEncodeTime
	buf  []byte
	used int

	encodeMessageSizeFn func(int)
	encodeMessageFn     func(metricpb.ChunkedMetricWithStoragePolicyWithEncodeTime) error
}

// NewAggregatedEncoder creates a new aggregated encoder.
func NewAggregatedEncoder(opts UnaggregatedOptions) AggregatedEncoder {
	e := &aggregatedEncoder{
		pool:           opts.BytesPool(),
		initBufSize:    opts.InitBufferSize(),
		maxMessageSize: opts.MaxMessageSize(),
	}
	e.encodeMessageSizeFn = e.encodeMessageSize
	e.encodeMessageFn = e.encodeMessage
	e.Reset(nil)
	return e
}

func (enc *aggregatedEncoder) Len() int { return enc.used }

func (enc *aggregatedEncoder) Reset(initData []byte) {
	if enc.buf != nil && enc.pool != nil {
		enc.pool.Put(enc.buf)
	}
	bufSize := int(math.Max(float64(enc.initBufSize), float64(len(initData))))
	enc.buf = allocate(enc.pool, bufSize)
	copy(enc.buf, initData)
	enc.used = len(initData)
}

func (enc *aggregatedEncoder) Truncate(n int) error {
	if n < 0 || n > enc.used {
		return fmt.Errorf("truncation out of range: used=%d, target=%d", enc.used, n)
	}
	enc.used = n
	return nil
}

func (enc *aggregatedEncoder) Relinquish() Buffer {
	res := NewBuffer(enc.buf[:enc.used], enc.pool)
	enc.buf = nil
	enc.used = 0
	return res
}

func (enc *aggregatedEncoder) EncodeChunkedMetricWithStoragePolicyAndEncodeTime(
	cmp aggregated.ChunkedMetricWithStoragePolicy,
	encodedAtNanos int64,
) error {
	if err := cmp.ToProto(&enc.cm.ChunkedMetricWithStoragePolicy); err != nil {
		return fmt.Errorf("counter with metadatas proto conversion failed: %v", err)
	}
	enc.cm.EncodeNanos = encodedAtNanos
	return enc.encode()
}

func (enc *aggregatedEncoder) encode() error {
	msgSize := enc.cm.Size()
	if msgSize > enc.maxMessageSize {
		return fmt.Errorf("message size %d is larger than maximum supported size %d", msgSize, enc.maxMessageSize)
	}
	payloadSize := maxMessageSizeInBytes + msgSize
	enc.ensureBufferSize(enc.used + payloadSize)
	enc.encodeMessageSizeFn(msgSize)
	return enc.encodeMessageFn(enc.cm)
}

// TODO(cw) Reuse
// ensureBufferSize ensures the internal buffer has at least the specified target size.
func (enc *aggregatedEncoder) ensureBufferSize(targetSize int) {
	enc.buf = ensureBufferSize(enc.buf, enc.pool, targetSize, copyData)
}

// TODO(cw) Reuse
func (enc *aggregatedEncoder) encodeMessageSize(msgSize int) {
	n := binary.PutVarint(enc.buf[enc.used:], int64(msgSize))
	enc.used += n
}

func (enc *aggregatedEncoder) encodeMessage(pb metricpb.ChunkedMetricWithStoragePolicyWithEncodeTime) error {
	n, err := pb.MarshalTo(enc.buf[enc.used:])
	if err != nil {
		return err
	}
	enc.used += n
	return nil
}
