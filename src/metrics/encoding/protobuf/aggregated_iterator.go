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
	"github.com/m3db/m3/src/metrics/metric/aggregated"
	"github.com/m3db/m3/src/metrics/policy"
)

// AggregatedIterator decodes aggregated metrics.
type AggregatedIterator interface {
	// Next returns true if there are more items to decode.
	Next() bool

	// Current returns the current decoded value.
	Current() (aggregated.RawMetric, policy.StoragePolicy, int64)

	// Err returns the error encountered during decoding, if any.
	Err() error

	// Close closes the iterator.
	Close()
}

// type aggregatedIterator struct {
// 	reader         encoding.ByteReadScanner
// 	bytesPool      pool.BytesPool
// 	maxMessageSize int

// 	closed  bool
// 	pb      metricpb.ChunkedMetricWithStoragePolicyWithEncodeTime
// 	chunked aggregated.ChunkedMetricWithStoragePolicy
// 	encoded int64
// 	sp      policy.StoragePolicy
// 	raw     aggregated.RawMetric
// 	buf     []byte
// 	err     error
// }

// // NewAggregatedIterator creates a new aggregated iterator.
// func NewAggregatedIterator(
// 	reader encoding.ByteReadScanner,
// 	opts UnaggregatedOptions,
// ) UnaggregatedIterator {
// 	bytesPool := opts.BytesPool()
// 	return &aggregatedIterator{
// 		reader:         reader,
// 		bytesPool:      bytesPool,
// 		maxMessageSize: opts.MaxMessageSize(),
// 		buf:            allocate(bytesPool, opts.InitBufferSize()),
// 	}
// }

// func (it *aggregatedIterator) Close() {
// 	if it.closed {
// 		return
// 	}
// 	it.closed = true
// 	it.reader = nil
// 	// it.pb.Reset()
// 	// it.msg = encoding.UnaggregatedMessageUnion{}
// 	if it.bytesPool != nil && it.buf != nil {
// 		it.bytesPool.Put(it.buf)
// 	}
// 	it.bytesPool = nil
// 	it.buf = nil
// 	it.err = nil
// }

// func (it *aggregatedIterator) Err() error { return it.err }
// func (it *aggregatedIterator) Current() (aggregated.RawMetric, policy.StoragePolicy, int64) {
// 	return it.raw, it.sp, it.encoded
// }

// func (it *aggregatedIterator) Next() bool {
// 	if it.err != nil || it.closed {
// 		return false
// 	}
// 	size, err := it.decodeSize()
// 	if err != nil {
// 		return false
// 	}
// 	if size > it.maxMessageSize {
// 		it.err = fmt.Errorf("decoded message size %d is larger than supported max message size %d", size, it.maxMessageSize)
// 		return false
// 	}
// 	it.ensureBufferSize(size)
// 	if err := it.decodeMessage(size); err != nil {
// 		return false
// 	}
// 	return true
// }

// func (it *aggregatedIterator) decodeSize() (int, error) {
// 	n, err := binary.ReadVarint(it.reader)
// 	if err != nil {
// 		it.err = err
// 		return 0, err
// 	}
// 	return int(n), nil
// }

// func (it *aggregatedIterator) ensureBufferSize(targetSize int) {
// 	it.buf = ensureBufferSize(it.buf, it.bytesPool, targetSize, dontCopyData)
// }

// func (it *aggregatedIterator) decodeMessage(size int) error {
// 	_, err := io.ReadFull(it.reader, it.buf[:size])
// 	if err != nil {
// 		it.err = err
// 		return err
// 	}
// 	resetChunkedMetricWithStoragePolicyWithEncodeNanosProto(&it.pb)
// 	if err := it.pb.Unmarshal(it.buf[:size]); err != nil {
// 		it.err = err
// 		return err
// 	}

// 	if err := it.chunked.FromProto(it.pb.ChunkedMetricWithStoragePolicy); err != nil {
// 		it.err = err
// 		return err
// 	}
// 	it.chunk
// 	return it.err
// }
