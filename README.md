Benchmark of my [maintained fork](https://github.com/joeycumines/statsd) of `alexcesaro/statsd`, with the
Go clients listed on the
[StatsD wiki](https://github.com/etsy/statsd/wiki#client-implementations):

```
$ go test -bench . -benchmem -benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/joeycumines/statsdbench
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkAlexcesaro/upstream-16         	 3653144	      1600 ns/op	        64.04 bytes/op	         0.0000005 emptyPackets/op	         0 invalidMetrics/op	         2.345 metrics/op	         0.04510 packets/op	         1.789 unexpectedMetrics/op	         2.345 validMetrics/op	     356 B/op	       0 allocs/op
BenchmarkAlexcesaro/fork-16             	 3599080	      1657 ns/op	        70.55 bytes/op	         0.0000006 emptyPackets/op	         0 invalidMetrics/op	         2.583 metrics/op	         0.04968 packets/op	         2.281 unexpectedMetrics/op	         2.583 validMetrics/op	     361 B/op	       0 allocs/op
BenchmarkAlexcesaro/fork_smaller_packets-16         	 1274977	      4760 ns/op	        81.17 bytes/op	         0.0000016 emptyPackets/op	         0 invalidMetrics/op	         2.976 metrics/op	         0.1654 packets/op	         0 unexpectedMetrics/op	         2.976 validMetrics/op	     422 B/op	       0 allocs/op
BenchmarkAlexcesaro/inline_flush-16                 	   77793	     76783 ns/op	        79.00 bytes/op	         0.0000257 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     583 B/op	       0 allocs/op
BenchmarkAlexcesaro/mem_server_buffered-16          	33485823	       371.5 ns/op	        82.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         0.05769 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     180 B/op	       0 allocs/op
BenchmarkAlexcesaro/mem_server_unbuffered-16        	29101744	       199.7 ns/op	        82.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     349 B/op	       0 allocs/op
BenchmarkCactus/buffered-16                         	 3322560	      1815 ns/op	        78.48 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         2.873 metrics/op	         0.05527 packets/op	         1.893 unexpectedMetrics/op	         2.873 validMetrics/op	     492 B/op	       0 allocs/op
BenchmarkCactus/buffered_smaller_packets-16         	 1000000	      5410 ns/op	        81.06 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         2.972 metrics/op	         0.1651 packets/op	         0 unexpectedMetrics/op	         2.972 validMetrics/op	     441 B/op	       0 allocs/op
BenchmarkCactus/unbuffered-16                       	   74956	     81110 ns/op	        79.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     797 B/op	       3 allocs/op
BenchmarkCactus/net_conn-16                         	   75384	     78940 ns/op	        79.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     602 B/op	       0 allocs/op
BenchmarkCactus/mem_server-16                       	26034400	       238.5 ns/op	        79.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     357 B/op	       0 allocs/op
BenchmarkG2s-16                                     	   73782	     83058 ns/op	        79.00 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	    1151 B/op	      19 allocs/op
BenchmarkQuipo-16                                   	   73635	     81897 ns/op	        78.99 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     776 B/op	       9 allocs/op
BenchmarkSmira/strange_results-16                   	38341719	       142.6 ns/op	         6.741 bytes/op	         0 emptyPackets/op	         0 invalidMetrics/op	         0.2468 metrics/op	         0.004746 packets/op	         0.1670 unexpectedMetrics/op	         0.2468 validMetrics/op	     175 B/op	       0 allocs/op
PASS
ok  	github.com/joeycumines/statsdbench	125.685s
```
