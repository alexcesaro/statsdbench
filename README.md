Benchmark of my [StatsD client](https://github.com/alexcesaro/statsd) with the
Go clients listed on the
[StatsD wiki](https://github.com/etsy/statsd/wiki#client-implementations):

```
$ go test -bench . -benchmem -benchtime=5s
BenchmarkAlexcesaro-4    5000000     1277 ns/op      0 B/op     0 allocs/op
BenchmarkCactus-4        2000000     3301 ns/op      4 B/op     0 allocs/op
BenchmarkG2s-4            200000    38715 ns/op    624 B/op    26 allocs/op
BenchmarkQuipo-4         2000000     4887 ns/op    496 B/op    10 allocs/op
```
