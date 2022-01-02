Benchmark of my [maintained fork](https://github.com/joeycumines/statsd) of `alexcesaro/statsd`, with the
Go clients listed on the
[StatsD wiki](https://github.com/etsy/statsd/wiki#client-implementations):

```
$ go test -bench . -benchmem -benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/joeycumines/statsdbench
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkAlexcesaro/upstream-16                 	 3694010	  1627 ns/op	    64.04 bytes/op	     0.0000005 emptyPackets/op	         0 invalidMetrics/op	     2.345 metrics/op	         0.04510 packets/op	         1.062 unexpectedMetrics/op	     2.345 validMetrics/op	     352 B/op	       0 allocs/op
BenchmarkAlexcesaro/fork-16                     	 3688959	  1634 ns/op	    70.14 bytes/op	     0.0000005 emptyPackets/op	         0 invalidMetrics/op	     2.568 metrics/op	         0.04940 packets/op	         2.104 unexpectedMetrics/op	     2.568 validMetrics/op	     440 B/op	       0 allocs/op
BenchmarkAlexcesaro/fork_smaller_packets-16     	 1250258	  4814 ns/op	    80.72 bytes/op	     0.0000016 emptyPackets/op	         0 invalidMetrics/op	     2.959 metrics/op	         0.1645 packets/op	         0 unexpectedMetrics/op	         2.959 validMetrics/op	     430 B/op	       0 allocs/op
BenchmarkAlexcesaro/inline_flush-16             	   79965	 76419 ns/op	    79.00 bytes/op	     0.0000250 emptyPackets/op	         0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     597 B/op	       0 allocs/op
BenchmarkAlexcesaro/mem_server_buffered-16      	35463738	 326.2 ns/op	    82.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         0.05769 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     170 B/op	       0 allocs/op
BenchmarkAlexcesaro/mem_server_unbuffered-16    	30631983	 194.5 ns/op	    82.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     332 B/op	       0 allocs/op
BenchmarkCactus/buffered-16                     	 3371697	  1740 ns/op	    76.75 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     2.810 metrics/op	         0.05405 packets/op	         1.939 unexpectedMetrics/op	     2.810 validMetrics/op	     485 B/op	       0 allocs/op
BenchmarkCactus/buffered_smaller_packets-16     	 1000000	  5232 ns/op	    81.19 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     2.976 metrics/op	         0.1654 packets/op	         0 unexpectedMetrics/op	         2.976 validMetrics/op	     441 B/op	       0 allocs/op
BenchmarkCactus/unbuffered-16                   	   73893	 77170 ns/op	    79.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     806 B/op	       3 allocs/op
BenchmarkCactus/net_conn-16                     	   80132	 75776 ns/op	    79.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     596 B/op	       0 allocs/op
BenchmarkCactus/mem_server-16                   	28030728	 229.9 ns/op	    79.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     360 B/op	       0 allocs/op
BenchmarkG2s-16                                 	   74582	 80919 ns/op	    79.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	    1144 B/op	      19 allocs/op
BenchmarkQuipo-16                               	   72733	 79852 ns/op	    79.00 bytes/op	     0 emptyPackets/op	                 0 invalidMetrics/op	     3.000 metrics/op	         3.000 packets/op	         0 unexpectedMetrics/op	         3.000 validMetrics/op	     684 B/op	       9 allocs/op
PASS
ok  	github.com/joeycumines/statsdbench	120.332s
```
