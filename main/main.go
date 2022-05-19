package main

import (
	"fmt"

	"github.com/shunsukuda/bytesmemorypool"
)

var mp = bytesmemorypool.NewMemoryPool()

func main() {
	b := mp.Get(0)
	fmt.Println(cap(b))
}
