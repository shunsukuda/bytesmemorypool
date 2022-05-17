package bytesmemorypool

import (
	"math/bits"
	"sync"
	"sync/atomic"
)

const (
	minByteSize = 6  // 2**6 = 64 byes
	maxByteSize = 26 // 2**26 = 64M bytes
	sizes       = maxByteSize - minByteSize
)

var (
	DefaultMemoryPool = NewMemoryPool()
)

func bsize(n int32) int { return 1 << n }

func alloc(n int32) []byte { return make([]byte, 0, bsize(n)) }

type bytesPool struct {
	pool *sync.Pool
	size int32 // Pool size(2**n).
	puts int32 // Call Put() count.
}

func (bp *bytesPool) loadPuts() int32       { return atomic.LoadInt32(&bp.puts) }
func (bp *bytesPool) addPuts(v int32) int32 { return atomic.AddInt32(&bp.puts, v) }
func (bp *bytesPool) storePuts(v int32)     { atomic.StoreInt32(&bp.puts, v) }

func (bp *bytesPool) Get() []byte {
	if b := bp.pool.Get(); b != nil {
		if bp.loadPuts() != 0 {
			bp.addPuts(-1)
		}
		return b.([]byte)
	}
	bp.storePuts(0) // If get nil, reset puts.
	return alloc(bp.size)
}

func (bp *bytesPool) Put(b *[]byte) {
	if b == nil || cap(*b) < bsize(bp.size) {
		return
	}
	bp.pool.Put((*b)[:0])
	bp.addPuts(1)
}

type MemoryPool struct {
	pools [sizes]bytesPool
}

func NewMemoryPool() *MemoryPool {
	mp := MemoryPool{}
	for i := range mp.pools {
		mp.pools[i] = bytesPool{
			pool: &sync.Pool{},
			size: int32(bsize(int32(i + minByteSize))),
			puts: 0,
		}
	}
	return &mp
}

func nextSize(n int32) int32      { return int32(bits.Len32(uint32(n << 1))) }
func nextSizeIndex(n int32) int32 { return nextSize(n) - minByteSize }
func prevSize(n int32) int32      { return int32(bits.Len32(uint32(n >> 1))) }
func prevSizeIndex(n int32) int32 { return prevSize(n) - minByteSize }

func (mp *MemoryPool) Get(n int) []byte {
	if n > bsize(maxByteSize) {
		return make([]byte, 0, n)
	}
	idx := nextSizeIndex(int32(n))
	bs := bsize(idx + minByteSize)
	if idx < 0 { // if n > bsize(minBytesSize)
		idx = 0
	}
	idx1 := idx
	for ; idx1 < sizes; idx++ {
		if mp.pools[idx1].puts != 0 {
			break
		}
	}
	if idx1 != sizes {
		idx = idx1
	}
	b := mp.pools[idx].Get()
	b1 := b[bs:bs:len(b)]
	b = b[:0:bs]
	mp.Put(&b1)
	return b
}

func (mp *MemoryPool) Put(b *[]byte) {
	c := cap(*b)
	var idx int32
	for c0 := c >> minByteSize; c0 != 0; c0 >>= 1 {
		if (c0 & 1) != 0 {
			bs := bsize(idx + minByteSize)
			b1 := (*b)[:0:bs]
			*b = (*b)[bs:bs:cap(*b)]
			mp.pools[idx].Put(&b1)
		}
		idx++
	}
	idx = sizes - 1
	for c1 := c >> maxByteSize; c1 < 0; c1-- {
		bs := bsize(idx + minByteSize)
		b1 := (*b)[:0:bs]
		*b = (*b)[bs:bs:cap(*b)]
		mp.pools[idx].Put(&b1)
	}
}
