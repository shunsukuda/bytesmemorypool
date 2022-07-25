package bytesmemorypool

import (
	"reflect"
	"sync"
	"testing"
)

func Test_bytesPool_Get(t *testing.T) {
	tests := []struct {
		name     string
		bp       *bytesPool
		doPuts   int
		wantPuts int32
		wantCap  int
	}{
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, doPuts: 0, wantPuts: 0, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 10, puts: 0}, doPuts: 0, wantPuts: 0, wantCap: 1024},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, doPuts: 5, wantPuts: 4, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 5}, doPuts: 0, wantPuts: 0, wantCap: 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.doPuts
			for i > 0 {
				b := make([]byte, 0, tt.wantCap)
				tt.bp.Put(&b)
				i--
			}
			got := tt.bp.Get()
			if puts := tt.bp.loadPuts(); puts != tt.wantPuts {
				t.Errorf("bytesPool.Get() bp = %v, puts = %v, want %v", tt.bp, puts, tt.wantPuts)
			}
			if c := cap(got); c != tt.wantCap {
				t.Errorf("bytesPool.Get() bp = %v, cap = %v, want %v", tt.bp, c, tt.wantCap)
			}
		})
		//runtime.GC()
	}
}

func Test_bytesPool_Put(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		bp       *bytesPool
		args     args
		wantPuts int32
		wantCap  int
	}{
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, args: args{b: nil}, wantPuts: 0, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, args: args{b: make([]byte, 0, 63)}, wantPuts: 0, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, args: args{b: make([]byte, 0, 64)}, wantPuts: 1, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 10, puts: 0}, args: args{b: make([]byte, 0, 1024)}, wantPuts: 1, wantCap: 1024},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, args: args{b: make([]byte, 0, 64)}, wantPuts: 1, wantCap: 64},
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 0}, args: args{b: make([]byte, 0, 64)}, wantPuts: 1, wantCap: 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bp.Put(&tt.args.b)
			//got := tt.bp.Get()
			if puts := tt.bp.loadPuts(); puts != tt.wantPuts {
				t.Errorf("bytesPool.Put() bp = %v, puts = %v, want %v", tt.bp, puts, tt.wantPuts)
			}
			/*
				if c := cap(got); c != tt.wantCap {
					t.Errorf("bytesPool.Put() bp = %v, cap = %v, want %v", tt.bp, c, tt.wantCap)
				}
			*/
		})
		//runtime.GC()
	}
}

func Test_nextSizeIndex(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{name: "", args: args{n: 0}, want: 0},
		{name: "", args: args{n: 64}, want: 0},
		{name: "", args: args{n: 100}, want: 1},
		{name: "", args: args{n: 1024}, want: 4},
		{name: "", args: args{n: 1025}, want: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextSizeIndex(tt.args.n); got != tt.want {
				t.Errorf("nextSizeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryPool_Get(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		mp      *MemoryPool
		args    args
		putSize int
		wantCap int
	}{
		{name: "", mp: NewMemoryPool(), args: args{n: 0}, putSize: 0, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 32}, putSize: 0, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 63}, putSize: 0, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 64}, putSize: 0, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 100}, putSize: 0, wantCap: 128},
		{name: "", mp: NewMemoryPool(), args: args{n: 1024}, putSize: 0, wantCap: 1024},
		{name: "", mp: NewMemoryPool(), args: args{n: 1025}, putSize: 0, wantCap: 2048},
		{name: "", mp: NewMemoryPool(), args: args{n: 0}, putSize: 64, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 0}, putSize: 128, wantCap: 64},
		{name: "", mp: NewMemoryPool(), args: args{n: 100}, putSize: 128, wantCap: 128},
		{name: "", mp: NewMemoryPool(), args: args{n: bsize(27) + 1}, putSize: 0, wantCap: bsize(27) + 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.putSize != 0 {
				b := make([]byte, 0, tt.putSize)
				tt.mp.Put(&b)
			}
			got := tt.mp.Get(tt.args.n)
			if c := cap(got); c != tt.wantCap {
				t.Errorf("MemoryPool.Get() cap = %v, want %v", c, tt.wantCap)
			}
		})
	}
}

func TestMemoryPool_Put(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name     string
		mp       *MemoryPool
		args     args
		wantPuts [21]int32
	}{
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 0)}, wantPuts: [21]int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 63)}, wantPuts: [21]int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 64)}, wantPuts: [21]int32{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 100)}, wantPuts: [21]int32{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 128)}, wantPuts: [21]int32{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 500)}, wantPuts: [21]int32{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 1000)}, wantPuts: [21]int32{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, bsize(27))}, wantPuts: [21]int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, bsize(28))}, wantPuts: [21]int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}},
		{name: "", mp: NewMemoryPool(), args: args{b: make([]byte, 0, 1000+bsize(28))}, wantPuts: [21]int32{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mp.Put(&tt.args.b)
			p := tt.mp.loadPuts()
			if !reflect.DeepEqual(p, tt.wantPuts) {
				t.Errorf("MemoryPool.Put() puts = %v, want %v", p, tt.wantPuts)
			}
		})
	}
}

func TestMakeByteSlice(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		wantCap int
	}{
		{name: "", args: args{n: 0}, wantCap: 64},
		{name: "", args: args{n: 32}, wantCap: 64},
		{name: "", args: args{n: 63}, wantCap: 64},
		{name: "", args: args{n: 64}, wantCap: 64},
		{name: "", args: args{n: 100}, wantCap: 128},
		{name: "", args: args{n: 1024}, wantCap: 1024},
		{name: "", args: args{n: 1025}, wantCap: 2048},
		{name: "", args: args{n: 0}, wantCap: 64},
		{name: "", args: args{n: 0}, wantCap: 64},
		{name: "", args: args{n: 100}, wantCap: 128},
		{name: "", args: args{n: bsize(27) + 1}, wantCap: bsize(27) + 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeByteSlice(tt.args.n); cap(got) != tt.wantCap {
				t.Errorf("MakeByteSlice() cap = %v, want %v", got, tt.wantCap)
			}
		})
	}
}
