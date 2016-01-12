package statsdbench

import (
	"net"
	"strconv"
	"testing"
	"time"

	alexcesaro "github.com/alexcesaro/statsd"
	cactus "github.com/cactus/go-statsd-client/statsd"
	"github.com/peterbourgon/g2s"
	quipo "github.com/quipo/statsd"
)

const (
	addr        = ":8125"
	prefix      = "prefix."
	prefixNoDot = "prefix"
	counterKey  = "foo.bar.counter"
	gaugeKey    = "foo.bar.gauge"
	gaugeValue  = 42
	timingKey   = "foo.bar.timing"
	tValDur     = 153 * time.Millisecond
	tValInt     = int(153)
	tValInt64   = int64(153)
	flushPeriod = 100 * time.Millisecond
)

type logger struct{}

func (logger) Println(v ...interface{}) {}

func BenchmarkAlexcesaro(b *testing.B) {
	s := newServer()
	c, err := alexcesaro.New(addr, alexcesaro.WithPrefix(prefix), alexcesaro.WithFlushPeriod(flushPeriod))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		c.Increment(counterKey)
		c.Gauge(gaugeKey, gaugeValue)
		c.Timing(timingKey, tValInt, 1)
	}
	c.Close()
	s.Close()
}

func BenchmarkCactus(b *testing.B) {
	s := newServer()
	c, err := cactus.NewBufferedClient(addr, prefix, flushPeriod, 1432)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		c.Inc(counterKey, 1, 1)
		c.Gauge(gaugeKey, gaugeValue, 1)
		c.Timing(timingKey, tValInt64, 1)
	}
	c.Close()
	s.Close()
}

func BenchmarkCactusTimingAsDuration(b *testing.B) {
	s := newServer()
	c, err := cactus.NewBufferedClient(addr, prefix, flushPeriod, 1432)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		c.Inc(counterKey, 1, 1)
		c.Gauge(gaugeKey, gaugeValue, 1)
		c.TimingDuration(timingKey, tValDur, 1)
	}
	c.Close()
	s.Close()
}

func BenchmarkG2s(b *testing.B) {
	s := newServer()
	c, err := g2s.Dial("udp", addr)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		c.Counter(1, counterKey, 1)
		c.Gauge(1, gaugeKey, strconv.Itoa(gaugeValue))
		c.Timing(1, timingKey, tValDur)
	}
	s.Close()
}

func BenchmarkQuipo(b *testing.B) {
	s := newServer()
	c := quipo.NewStatsdBuffer(flushPeriod, quipo.NewStatsdClient(addr, prefix))
	c.Logger = logger{}
	for i := 0; i < b.N; i++ {
		c.Incr(counterKey, 1)
		c.Gauge(gaugeKey, gaugeValue)
		c.Timing(timingKey, tValInt64)
	}
	c.Close()
	s.Close()
}

func BenchmarkQuipoTimingAsDuration(b *testing.B) {
	s := newServer()
	c := quipo.NewStatsdBuffer(flushPeriod, quipo.NewStatsdClient(addr, prefix))
	c.Logger = logger{}
	for i := 0; i < b.N; i++ {
		c.Incr(counterKey, 1)
		c.Gauge(gaugeKey, gaugeValue)
		c.PrecisionTiming(timingKey, tValDur)
	}
	c.Close()
	s.Close()
}

type server struct {
	conn   *net.UDPConn
	closed chan bool
}

func newServer() *server {
	addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	s := &server{conn: conn, closed: make(chan bool)}
	go func() {
		buf := make([]byte, 512)
		for {
			_, err := conn.Read(buf)
			if err != nil {
				s.closed <- true
				return
			}
		}
	}()
	return s
}

func (s *server) Close() {
	s.conn.Close()
	<-s.closed
	return
}
