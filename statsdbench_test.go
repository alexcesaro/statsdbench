package statsdbench

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	cactus "github.com/cactus/go-statsd-client/v5/statsd"
	ac_fork "github.com/joeycumines/statsd"
	"github.com/peterbourgon/g2s"
	quipo "github.com/quipo/statsd"
	smira "github.com/smira/go-statsd"
	ac_upstream "gopkg.in/alexcesaro/statsd.v2"
)

const (
	addr               = "127.0.0.1:0"
	prefix             = "prefix."
	prefixNoDot        = "prefix"
	counterKey         = "foo.bar.counter"
	gaugeKey           = "foo.bar.gauge"
	gaugeValue         = 42
	timingKey          = "foo.bar.timing"
	timingValue        = 153 * time.Millisecond
	flushPeriod        = 100 * time.Millisecond
	maxPacketSizeUpper = 1432
	maxPacketSizeLower = 508
)

type (
	Server interface {
		Close()
		Buffer() []byte
		Packets() []int
		Addr() string
		Writer() io.WriteCloser
	}

	ServerFactory func() Server

	udpServer struct {
		conn       *net.UDPConn
		clientConn net.Conn
		closed     chan struct{}
		buffer     []byte
		packets    []int
	}

	memServer struct {
		buffer  bytes.Buffer
		packets []int
	}

	memServerWriter struct {
		srv *memServer
	}

	cactusSender struct {
		io.WriteCloser
		mu sync.Mutex
	}

	logger struct{}
)

func BenchmarkAlexcesaro(b *testing.B) {
	type Client interface {
		Increment(bucket string)
		Gauge(bucket string, value interface{})
		Timing(bucket string, value interface{})
		Close()
	}
	for _, tc := range [...]struct {
		Name   string
		Server ServerFactory
		Client func(server Server) (Client, error)
	}{
		{
			Name:   `upstream`, // no longer maintained
			Server: newUDPServer,
			Client: func(server Server) (Client, error) {
				return ac_upstream.New(
					ac_upstream.Address(server.Addr()),
					ac_upstream.Prefix(prefixNoDot),
					ac_upstream.FlushPeriod(flushPeriod),
					ac_upstream.MaxPacketSize(maxPacketSizeUpper),
				)
			},
		},
		{
			Name:   `fork`, // backwards compatible, maintained
			Server: newUDPServer,
			Client: func(server Server) (Client, error) {
				return ac_fork.New(
					ac_fork.Address(server.Addr()),
					ac_fork.Prefix(prefixNoDot),
					ac_fork.FlushPeriod(flushPeriod),
					ac_fork.MaxPacketSize(maxPacketSizeUpper),
				)
			},
		},
		{
			Name:   `fork smaller packets`, // backwards compatible, maintained
			Server: newUDPServer,
			Client: func(server Server) (Client, error) {
				return ac_fork.New(
					ac_fork.Address(server.Addr()),
					ac_fork.Prefix(prefixNoDot),
					ac_fork.FlushPeriod(flushPeriod),
					ac_fork.MaxPacketSize(maxPacketSizeLower),
				)
			},
		},
		{
			Name:   `inline flush`, // one packet per metric
			Server: newUDPServer,
			Client: func(server Server) (Client, error) {
				return ac_fork.New(
					ac_fork.Address(server.Addr()),
					ac_fork.Prefix(prefixNoDot),
					ac_fork.InlineFlush(true),
				)
			},
		},
		{
			Name:   `mem server buffered`,
			Server: newMemServer,
			Client: func(server Server) (Client, error) {
				return ac_fork.New(
					ac_fork.WriteCloser(server.Writer()),
					ac_fork.Prefix(prefixNoDot),
					ac_fork.FlushPeriod(flushPeriod),
					ac_fork.MaxPacketSize(maxPacketSizeUpper),
				)
			},
		},
		{
			Name:   `mem server unbuffered`,
			Server: newMemServer,
			Client: func(server Server) (Client, error) {
				return ac_fork.New(
					ac_fork.WriteCloser(server.Writer()),
					ac_fork.Prefix(prefixNoDot),
					ac_fork.InlineFlush(true),
				)
			},
		},
	} {
		b.Run(tc.Name, func(b *testing.B) {
			s := tc.Server()
			defer s.Close()

			c, err := tc.Client(s)
			if err != nil {
				b.Fatal(err)
			}
			defer c.Close()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				c.Increment(counterKey)
				c.Gauge(gaugeKey, gaugeValue)
				c.Timing(timingKey, int64(timingValue/time.Millisecond))
			}

			c.Close()

			b.StopTimer()

			s.Close()
			validateServer(b, s)
		})
	}
}

func BenchmarkCactus(b *testing.B) {
	for _, tc := range [...]struct {
		Name   string
		Server ServerFactory
		Client func(server Server) (cactus.Statter, error)
	}{
		{
			Name:   `buffered`,
			Server: newUDPServer,
			Client: func(server Server) (cactus.Statter, error) {
				return cactus.NewBufferedClient(server.Addr(), prefixNoDot, flushPeriod, maxPacketSizeUpper)
			},
		},
		{
			Name:   `buffered smaller packets`,
			Server: newUDPServer,
			Client: func(server Server) (cactus.Statter, error) {
				return cactus.NewBufferedClient(server.Addr(), prefixNoDot, flushPeriod, maxPacketSizeLower)
			},
		},
		{
			Name:   `unbuffered`,
			Server: newUDPServer,
			Client: func(server Server) (cactus.Statter, error) {
				return cactus.NewClientWithConfig(&cactus.ClientConfig{
					Address: server.Addr(),
					Prefix:  prefixNoDot,
				})
			},
		},
		{
			Name:   `net conn`,
			Server: newUDPServer,
			Client: func(server Server) (cactus.Statter, error) {
				return cactus.NewClientWithSender(&cactusSender{WriteCloser: server.Writer()}, prefixNoDot, 0)
			},
		},
		{
			Name:   `mem server`,
			Server: newMemServer,
			Client: func(server Server) (cactus.Statter, error) {
				return cactus.NewClientWithSender(&cactusSender{WriteCloser: server.Writer()}, prefixNoDot, 0)
			},
		},
	} {
		b.Run(tc.Name, func(b *testing.B) {
			s := tc.Server()
			defer s.Close()

			c, err := tc.Client(s)
			if err != nil {
				b.Fatal(err)
			}
			defer c.Close()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if err := c.Inc(counterKey, 1, 1); err != nil {
					b.Fatal(err)
				}
				if err := c.Gauge(gaugeKey, gaugeValue, 1); err != nil {
					b.Fatal(err)
				}
				if err := c.Timing(timingKey, int64(timingValue/time.Millisecond), 1); err != nil {
					b.Fatal(err)
				}
			}

			if err := c.Close(); err != nil {
				b.Fatal(err)
			}

			b.StopTimer()

			s.Close()
			validateServer(b, s)
		})
	}
}

func BenchmarkG2s(b *testing.B) {
	s := newUDPServer()
	defer s.Close()

	c, err := g2s.Dial("udp", s.Addr())
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Counter(1, prefix+counterKey, 1)
		c.Gauge(1, prefix+gaugeKey, strconv.Itoa(gaugeValue))
		c.Timing(1, prefix+timingKey, timingValue)
	}

	b.StopTimer()

	s.Close()
	validateServer(b, s)
}

func BenchmarkQuipo(b *testing.B) {
	s := newUDPServer()
	defer s.Close()

	// the buffered implementation seems to hang on close, probably a deadlock...
	// c := quipo.NewStatsdBuffer(flushPeriod, quipo.NewStatsdClient(s.Addr(), prefix))
	c := quipo.NewStatsdClient(s.Addr(), prefix)
	defer c.Close()
	c.Logger = logger{}

	if err := c.CreateSocket(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := c.Incr(counterKey, 1); err != nil {
			b.Fatal(err)
		}
		if err := c.Gauge(gaugeKey, gaugeValue); err != nil {
			b.Fatal(err)
		}
		if err := c.Timing(timingKey, int64(timingValue/time.Millisecond)); err != nil {
			b.Fatal(err)
		}
	}

	if err := c.Close(); err != nil {
		b.Fatal(err)
	}

	b.StopTimer()

	s.Close()
	validateServer(b, s)
}

func BenchmarkSmira(b *testing.B) {
	for _, tc := range [...]struct {
		Name   string
		Server ServerFactory
		Client func(server Server) (*smira.Client, error)
	}{
		{
			Name:   `strange results`,
			Server: newUDPServer,
			Client: func(server Server) (*smira.Client, error) {
				client := smira.NewClient(
					server.Addr(),
					smira.Logger(logger{}),
					smira.MaxPacketSize(maxPacketSizeUpper),
					smira.FlushInterval(flushPeriod),
					smira.MetricPrefix(prefix),
				)
				return client, nil
			},
		},
	} {
		b.Run(tc.Name, func(b *testing.B) {
			s := tc.Server()
			defer s.Close()

			c, err := tc.Client(s)
			if err != nil {
				b.Fatal(err)
			}
			defer c.Close()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				c.Incr(counterKey, 1)
				c.Gauge(gaugeKey, gaugeValue)
				c.Timing(timingKey, int64(timingValue/time.Millisecond))
			}

			if err := c.Close(); err != nil {
				b.Fatal(err)
			}

			b.StopTimer()

			s.Close()
			validateServer(b, s)
		})
	}
}

func validateServer(b *testing.B, s Server) {
	const (
		numMetrics = 3
	)
	//if l := len(s.packets); l <= 0 || l > b.N*numMetrics {
	//	b.Error(l)
	//}
	var (
		offset            int
		emptyPackets      int
		validMetrics      int
		invalidMetrics    int
		unexpectedMetrics int
		buffer            = make([]string, 0, numMetrics-1)
		metricIndexes     = map[string]int{
			`prefix.foo.bar.counter:1|c`:   0,
			`prefix.foo.bar.gauge:42|g`:    1,
			`prefix.foo.bar.timing:153|ms`: 2,
		}
		validateMetric = func(metric string) (expected bool, ok bool) {
			var i int
			i, ok = metricIndexes[metric]
			if ok {
				expected = len(buffer) == i
			}
			if !ok || len(buffer) >= cap(buffer) {
				buffer = buffer[:0]
			} else {
				buffer = append(buffer, metric)
			}
			return
		}
	)
	for i, n := range s.Packets() {
		if i < 3 {
			//b.Logf("packet %d: %q", i, string(s.Buffer()[offset:offset+n]))
		}
		if n == 0 {
			emptyPackets++
		} else {
			for _, metric := range strings.Split(string(s.Buffer()[offset:offset+n]), "\n") {
				if metric == `` {
					continue
				}
				expected, ok := validateMetric(metric)
				if ok {
					validMetrics++
				} else {
					invalidMetrics++
					b.Logf(`invalid metric: %q %q`, buffer, metric)
				}
				if ok && !expected {
					unexpectedMetrics++
					//b.Logf(`unexpected metric: %q %q`, buffer, metric)
				}
			}
		}
		offset += n
	}
	b.ReportMetric(float64(len(s.Buffer()))/float64(b.N), `bytes/op`)
	b.ReportMetric(float64(len(s.Packets()))/float64(b.N), `packets/op`)
	b.ReportMetric(float64(emptyPackets)/float64(b.N), `emptyPackets/op`)
	b.ReportMetric(float64(validMetrics+invalidMetrics)/float64(b.N), `metrics/op`)
	b.ReportMetric(float64(validMetrics)/float64(b.N), `validMetrics/op`)
	b.ReportMetric(float64(invalidMetrics)/float64(b.N), `invalidMetrics/op`)
	b.ReportMetric(float64(unexpectedMetrics)/float64(b.N), `unexpectedMetrics/op`)
}

func newUDPServer() Server {
	addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	s := &udpServer{
		conn:   conn,
		closed: make(chan struct{}),
	}
	s.clientConn, err = net.Dial(`udp`, s.Addr())
	if err != nil {
		panic(err)
	}
	go func() {
		buf := make([]byte, 1<<12)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err.(*net.OpError).Err.Error() != `use of closed network connection` {
					panic(err)
				}
				_ = s.conn.Close()
				close(s.closed)
				return
			}
			s.buffer = append(s.buffer, buf[:n]...)
			s.packets = append(s.packets, n)
		}
	}()
	return s
}

func newMemServer() Server { return new(memServer) }

func (logger) Println(v ...interface{}) {}

func (logger) Printf(fmt string, args ...interface{}) {}

func (s *udpServer) Addr() string {
	return s.conn.LocalAddr().String()
}

func (s *udpServer) Close() {
	s.clientConn.Close()
	s.conn.Close()
	<-s.closed
	return
}

func (s *udpServer) Buffer() []byte { return s.buffer }

func (s *udpServer) Packets() []int { return s.packets }

func (s *udpServer) Writer() io.WriteCloser { return s.clientConn }

func (x *memServer) Addr() string { return `` }

func (x *memServer) Close() {}

func (x *memServer) Buffer() []byte { return x.buffer.Bytes() }

func (x *memServer) Packets() []int { return x.packets }

func (x *memServer) Writer() io.WriteCloser { return memServerWriter{x} }

func (x memServerWriter) Write(b []byte) (n int, err error) {
	n, err = x.srv.buffer.Write(b)
	x.srv.packets = append(x.srv.packets, n)
	return
}

func (x memServerWriter) Close() error { return nil }

func (x *cactusSender) Send(data []byte) (int, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.WriteCloser.Write(data)
}
