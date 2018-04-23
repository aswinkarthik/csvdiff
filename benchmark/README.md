## Benchmark Results

Benchmark test can be found [here](https://github.com/aswinkarthik93/csvdiff/blob/master/pkg/digest/digest_benchmark_test.go).

```bash
$ cd ./pkg/digest
$ go test -bench=. -v -benchmem -benchtime=5s -cover
```

|                              |            |                         |                      |                     |
| ---------------------------- | ---------- | ----------------------- | -------------------- | ------------------- |
| BenchmarkCreate1-8           |    2000000 |             5967 ns/op  |           5474 B/op  |        21 allocs/op |
| BenchmarkCreate10-8          |     500000 |            16251 ns/op  |          10889 B/op  |        94 allocs/op |
| BenchmarkCreate100-8         |     100000 |           114219 ns/op  |          67139 B/op  |       829 allocs/op |
| BenchmarkCreate1000-8        |      10000 |          1042723 ns/op  |         674239 B/op  |      8078 allocs/op |
| BenchmarkCreate10000-8       |       1000 |         10386850 ns/op  |        6533806 B/op  |     80306 allocs/op |
| BenchmarkCreate100000-8      |        100 |        108740944 ns/op  |       64206718 B/op  |    804208 allocs/op |
| BenchmarkCreate1000000-8     |          5 |       1161730558 ns/op  |       672048142 B/op |  8039026 allocs/op  |
| BenchmarkCreate10000000-8    |          1 |       12721982424 ns/op |       6549111872 B/op| 80308455 allocs/op  |
