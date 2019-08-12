package digest

import (
	"fmt"
	"io"
	"testing"
)

const SomeText = "something-name-%d,346345ty,fdhfdh,5436456,gfgjfgj,45234545,nfhgjfgj,45745745,djhgfjfgj"

func BenchmarkCreate1(b *testing.B)     { benchmarkCreate(1, b) }
func BenchmarkCreate10(b *testing.B)    { benchmarkCreate(10, b) }
func BenchmarkCreate100(b *testing.B)   { benchmarkCreate(100, b) }
func BenchmarkCreate1000(b *testing.B)  { benchmarkCreate(1000, b) }
func BenchmarkCreate10000(b *testing.B) { benchmarkCreate(10000, b) }

func BenchmarkCreate100000(b *testing.B)   { benchmarkCreate(100000, b) }
func BenchmarkCreate1000000(b *testing.B)  { benchmarkCreate(1000000, b) }
func BenchmarkCreate10000000(b *testing.B) { benchmarkCreate(10000000, b) }

func benchmarkCreate(limit int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateDigestFor(limit, b)
	}
}

func CreateDigestFor(count int, b *testing.B) {
	b.StopTimer()
	reader := &Reader{limit: count}

	config := &Config{
		Reader: reader,
		Key:    []int{0},
		Value:  []int{1},
	}

	b.StartTimer()
	_, _, _ = Create(config)
}

type Reader struct {
	counter int
	limit   int
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.counter == r.limit {
		return 0, io.EOF
	}
	toRead := fmt.Sprintf("%d,%s\n", r.counter, SomeText)
	r.counter++
	return copy(p, toRead), nil
}
