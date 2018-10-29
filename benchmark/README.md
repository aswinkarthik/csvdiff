# Comparison with other tools

## Setup

* Using the majestic million data. (Source in credits section)
* Both files have 998390 rows and 12 columns.
* Only one modification between both files.
* Ran on Processor: Intel Core i7 2.5 GHz 4 cores 16 GB RAM

## Baseline

0. csvdiff (this tool) : *0m1.159s*

```bash
 time csvdiff majestic_million.csv majestic_million_diff.csv
Additions 0
Modifications 1
...

real	0m1.159s
user	0m2.167s
sys		0m0.222s
```

## Other tools

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

Benchmark test can be found [here](https://github.com/aswinkarthik/csvdiff/blob/master/pkg/digest/digest_benchmark_test.go).

```bash
$ cd ./pkg/digest
$ go test -bench=. -v -benchmem -benchtime=5s -cover
```

```
BenchmarkCreate1-8          	  200000	     31794 ns/op	  116163 B/op	      24 allocs/op
BenchmarkCreate10-8         	  200000	     43351 ns/op	  119993 B/op	      79 allocs/op
BenchmarkCreate100-8        	   50000	    142645 ns/op	  160577 B/op	     634 allocs/op
BenchmarkCreate1000-8       	   10000	    907308 ns/op	  621694 B/op	    6085 allocs/op
BenchmarkCreate10000-8      	    1000	   7998083 ns/op	 5117977 B/op	   60345 allocs/op
BenchmarkCreate100000-8     	     100	  81260585 ns/op	49106849 B/op	  604563 allocs/op
BenchmarkCreate1000000-8    	      10	 788485738 ns/op	520115434 B/op	 6042650 allocs/op
BenchmarkCreate10000000-8   	       1	7878009695 ns/op	5029061632 B/op	60346535 allocs/op
```