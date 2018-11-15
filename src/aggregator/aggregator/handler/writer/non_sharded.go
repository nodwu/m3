// Copyright (c) 2017 Uber Technologies, Inc.
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

package writer

import (
	"errors"
	"math/rand"

	"github.com/m3db/m3/src/aggregator/sharding"
	"github.com/m3db/m3/src/metrics/encoding/msgpack"
	"github.com/m3db/m3/src/metrics/encoding/protobuf"
	"github.com/m3db/m3/src/metrics/metric/aggregated"
	"github.com/m3db/m3/src/msg/producer"
	"github.com/m3db/m3x/clock"

	"github.com/uber-go/tally"
)

var (
	errWriterClosed = errors.New("writer is closed")
)

type nonShardedWriterMetrics struct {
	writerClosed  tally.Counter
	encodeSuccess tally.Counter
	encodeErrors  tally.Counter
	routeSuccess  tally.Counter
	routeErrors   tally.Counter
}

func newNonShardedWriterMetrics(scope tally.Scope) nonShardedWriterMetrics {
	encodeScope := scope.SubScope("encode")
	routeScope := scope.SubScope("route")
	return nonShardedWriterMetrics{
		writerClosed:  scope.Counter("writer-closed"),
		encodeSuccess: encodeScope.Counter("success"),
		encodeErrors:  encodeScope.Counter("errors"),
		routeSuccess:  routeScope.Counter("success"),
		routeErrors:   routeScope.Counter("errors"),
	}
}

// nonShardedWriter encodes data in a shard-aware fashion and routes them to the backend.
// nonShardedWriter is not thread safe.
type nonShardedWriter struct {
	sharding.AggregatedSharder

	nowFn                    clock.NowFn
	maxBufferSize            int
	encodingTimeSamplingRate float64
	encoder                  protobuf.Encoder
	p                        producer.Producer

	closed  bool
	rand    *rand.Rand
	metrics nonShardedWriterMetrics
	shardFn shardFn
	randFn  randFn
}

// NewShardedWriter creates a new sharded writer.
// func NewShardedWriter(
// 	sharderID sharding.SharderID,
// 	opts Options,
// ) (Writer, error) {
// 	sharder, err := sharding.NewAggregatedSharder(sharderID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	numShards := sharderID.NumShards()
// 	nowFn := opts.ClockOptions().NowFn()
// 	instrumentOpts := opts.InstrumentOptions()
// 	w := &nonShardedWriter{
// 		AggregatedSharder:        sharder,
// 		nowFn:                    nowFn,
// 		maxBufferSize:            opts.MaxBufferSize(),
// 		encodingTimeSamplingRate: opts.EncodingTimeSamplingRate(),
// 		rand:                     rand.New(rand.NewSource(nowFn().UnixNano())),
// 		metrics:                  newNonShardedWriterMetrics(instrumentOpts.MetricsScope()),
// 	}
// 	w.shardFn = w.Shard
// 	w.randFn = w.rand.Float64
// 	return w, nil
// }

func (w *nonShardedWriter) Write(mp aggregated.ChunkedMetricWithStoragePolicy) error {
	if w.closed {
		w.metrics.writerClosed.Inc(1)
		return errWriterClosed
	}
	return w.encode(mp)
}

func (w *nonShardedWriter) Flush() error {
	return nil
	// if w.closed {
	// 	w.metrics.writerClosed.Inc(1)
	// 	return errWriterClosed
	// }
	// multiErr := xerrors.NewMultiError()
	// for shard, encoder := range w.encodersByShard {
	// 	if encoder == nil {
	// 		continue
	// 	}
	// 	bufferedEncoder := encoder.Encoder()
	// 	buffer := bufferedEncoder.Buffer()
	// 	if buffer.Len() == 0 {
	// 		continue
	// 	}
	// 	newBufferedEncoder := w.bufferedEncoderPool.Get()
	// 	newBufferedEncoder.Reset()
	// 	encoder.Reset(newBufferedEncoder)
	// 	if err := w.Route(uint32(shard), common.NewRefCountedBuffer(bufferedEncoder)); err != nil {
	// 		w.metrics.routeErrors.Inc(1)
	// 		multiErr = multiErr.Add(err)
	// 	} else {
	// 		w.metrics.routeSuccess.Inc(1)
	// 	}
	// }
	// return multiErr.FinalError()
}

func (w *nonShardedWriter) Close() error {
	if w.closed {
		w.metrics.writerClosed.Inc(1)
		return errWriterClosed
	}
	err := w.Flush()
	w.closed = true
	return err
}

func (w *nonShardedWriter) encode(
	mp aggregated.ChunkedMetricWithStoragePolicy,
) error {
	//Use this later
	// shard := w.shardFn(mp.ChunkedID)
	var shard uint32
	var encodeNanos int64
	if w.encodingTimeSamplingRate > 0 && w.randFn() < w.encodingTimeSamplingRate {
		encodeNanos = w.nowFn().UnixNano()
	}
	if err := w.encoder.EncodeChunkedMetricAndEncodeTime(mp, encodeNanos); err != nil {
		w.metrics.encodeErrors.Inc(1)
		return err
	}
	w.metrics.encodeSuccess.Inc(1)
	if err := w.p.Produce(newMessage(shard, w.encoder.Buffer())); err != nil {
		w.metrics.routeErrors.Inc(1)
	}
	w.metrics.routeSuccess.Inc(1)
	return nil
}

type message struct {
	shard uint32
	data  msgpack.Buffer
}

// TODO(cw): Pool the messages if needed.
func newMessage(shard uint32, data msgpack.Buffer) producer.Message {
	return message{shard: shard, data: data}
}

func (d message) Shard() uint32 {
	return d.shard
}

func (d message) Bytes() []byte {
	return d.data.Bytes()
}

func (d message) Size() int {
	// Use the cap of the underlying byte slice in the buffer instead of
	// the length of the byte encoded to avoid "memory leak", for example
	// when the underlying buffer is 2KB, and it only encoded 300B, if we
	// use 300 as the size, then a producer with a buffer of 3GB could be
	// actually buffering 20GB in total for the underlying buffers.
	return cap(d.data.Bytes())
}

func (d message) Finalize(producer.FinalizeReason) {
	d.data.Close()
}
