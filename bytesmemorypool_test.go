package bytesmemorypool

import (
	"reflect"
	"sync"
	"testing"
)

func Test_bsize(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "", args: args{n: 6}, want: 64},
		{name: "", args: args{n: 10}, want: 1024},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bsize(tt.args.n); got != tt.want {
				t.Errorf("bsize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesPool_loadPuts(t *testing.T) {
	tests := []struct {
		name string
		bp   *bytesPool
		want int32
	}{
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 10}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bp.loadPuts(); got != tt.want {
				t.Errorf("bytesPool.loadPuts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesPool_addPuts(t *testing.T) {
	type args struct {
		v int32
	}
	tests := []struct {
		name string
		bp   *bytesPool
		args args
		want int32
	}{
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 10}, args: args{v: 1}, want: 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bp.addPuts(tt.args.v); got != tt.want {
				t.Errorf("bytesPool.addPuts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesPool_storePuts(t *testing.T) {
	type args struct {
		v int32
	}
	tests := []struct {
		name string
		bp   *bytesPool
		args args
		want int32
	}{
		{name: "", bp: &bytesPool{pool: &sync.Pool{}, size: 6, puts: 10}, args: args{v: 100}, want: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bp.storePuts(tt.args.v)
			got := tt.bp.loadPuts()
			if got != tt.want {
				t.Errorf("bytesPool.storePuts() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestNewMemoryPool(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryPool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryPool(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextSize(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextSize(tt.args.n); got != tt.want {
				t.Errorf("nextSize() = %v, want %v", got, tt.want)
			}
		})
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextSizeIndex(tt.args.n); got != tt.want {
				t.Errorf("nextSizeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prevSize(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prevSize(tt.args.n); got != tt.want {
				t.Errorf("prevSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prevSizeIndex(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prevSizeIndex(tt.args.n); got != tt.want {
				t.Errorf("prevSizeIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryPool_Get(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		mp   *MemoryPool
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mp.Get(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemoryPool.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryPool_Put(t *testing.T) {
	type args struct {
		b *[]byte
	}
	tests := []struct {
		name string
		mp   *MemoryPool
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mp.Put(tt.args.b)
		})
	}
}

func Test_alloc(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alloc(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("alloc() = %v, want %v", got, tt.want)
			}
		})
	}
}
