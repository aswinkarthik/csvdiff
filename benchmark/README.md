## Comparison with other tools


### Setup

* Using the majestic million data. (Source in credits section)
* Both files have 998390 rows and 12 columns.
* Only one modification between both files.
* Ran on Processor: Intel Core i7 2.5 GHz 4 cores 16 GB RAM

0. csvdiff (this tool) : *0m2.085s*

```bash
time csvdiff run -b majestic_million.csv -d majestic_million_diff.csv

# Additions: 0
# Modifications: 1

real	0m2.085s
user	0m3.861s
sys	0m0.340s
```

1. [data.table](https://github.com/Rdatatable/data.table) : *0m4.284s*

	* Join both csvs using `id` column.
	* Check inequality between both columns
	* Rscript in [data-table.r](/benchmark/data-table.r) (Can it be written better? New to R)

```bash
time Rscript data-table.r

real	0m4.284s
user	0m3.887s
sys	0m0.284s
```

2. [csvdiff](https://pypi.org/project/csvdiff/) written in Python : *0m48.115s*

```bash
time csvdiff --style=summary id majestic_million.csv majestic_million_diff.csv
0 rows removed (0.0%)
0 rows added (0.0%)
1 rows changed (0.0%)

real	0m48.115s
user	0m42.895s
sys	0m3.948s
```

3. GNU diff (Fastest) : *0m0.297s*

	* Seems the fastest. Couldn't even come close here.
	* However, it does line by line diff. Does not support compound keys of a csv or selective compare of columns. Hence the disclaimer, cannot be used a generic diff tool.
	* On another note, lets see if we can reach this.

```bash
time diff majestic_million.csv majestic_million_diff.csv

real	0m0.297s
user	0m0.144s
sys	0m0.147s
```

## Go Benchmark Results

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
