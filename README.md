Benchmark of my [StatsD client](https://github.com/alexcesaro/statsd) with the
Go clients listed on the
[StatsD wiki](https://github.com/etsy/statsd/wiki#client-implementations):

```
$ go test -bench . -benchmem -benchtime=5s
BenchmarkAlexcesaro-8            	20000000	       393 ns/op	       0 B/op	       0 allocs/op
BenchmarkCactus-8                	 2000000	      3138 ns/op	      50 B/op	       3 allocs/op
BenchmarkCactusTimingAsDuration-8	 2000000	      3307 ns/op	      82 B/op	       4 allocs/op
BenchmarkDieterbe-8              	 1000000	     12746 ns/op	     352 B/op	      19 allocs/op
BenchmarkG2s-8                   	 1000000	     12251 ns/op	     624 B/op	      26 allocs/op
BenchmarkQuipo-8                 	 5000000	      1852 ns/op	     400 B/op	       7 allocs/op
BenchmarkQuipoTimingAsDuration-8 	 5000000	      1437 ns/op	     192 B/op	       6 allocs/op
```
