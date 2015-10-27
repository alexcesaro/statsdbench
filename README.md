Benchmark of my [StatsD client](https://github.com/alexcesaro/statsd) with the
Go clients listed on the
[StatsD wiki](https://github.com/etsy/statsd/wiki#client-implementations):

```
$ go test -bench . -benchmem -benchtime=5s
BenchmarkAlexcesaro-4	10000000	    691 ns/op	     0 B/op	     0 allocs/op
BenchmarkCactus-4    	 1000000	   6845 ns/op	   164 B/op	     6 allocs/op
BenchmarkG2s-4       	  500000	  17032 ns/op	   624 B/op	    26 allocs/op
BenchmarkQuipo-4     	 3000000	   2818 ns/op	   400 B/op	     7 allocs/op
```
