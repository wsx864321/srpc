package pool

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"testing"
	"time"
)

var (
	InitialCap = 5
	MaxIdleCap = 10
	MaximumCap = 100
	network    = "tcp"
	address    = "127.0.0.1:7777"
	//factory    = func() (interface{}, error) { return net.Dial(network, address) }
	factory = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return net.Dial("tcp", address)
	}
	closeFac = func(v net.Conn) error {
		return v.Close()
	}
)

func init() {
	// used for factory function
	go rpcServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled

	rand.Seed(time.Now().UTC().UnixNano())
}

func TestNew(t *testing.T) {
	p, err := newChannelPool()
	defer p.release()
	if err != nil {
		t.Errorf("New error: %s", err)
	}
}
func TestPool_Get_Impl(t *testing.T) {
	p, _ := newChannelPool()
	defer p.release()

	conn, err := p.get(context.Background(), network, address)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	p.put(conn)
}

func TestPool_Get(t *testing.T) {
	p, _ := newChannelPool()
	defer p.release()

	_, err := p.get(context.Background(), network, address)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	// after one get, current capacity should be lowered by one.
	if p.len() != (InitialCap - 1) {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.len())
	}

	// get them all
	var wg sync.WaitGroup
	for i := 0; i < (MaximumCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.get(context.Background(), network, address)
			if err != nil {
				t.Errorf("Get error: %s", err)
			}
		}()
	}
	wg.Wait()

	if p.len() != 0 {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.len())
	}

	_, err = p.get(context.Background(), network, address)
	if err != ErrMaxActiveConnReached {
		t.Errorf("Get error: %s", err)
	}

}

func TestPool_Put(t *testing.T) {
	p, err := NewChannelPool(WithInitialCap(InitialCap), WithMaxCap(MaximumCap), WithFactory(factory), WithClose(closeFac), WithIdleTimeout(time.Second*20))
	if err != nil {
		t.Fatal(err)
	}
	defer p.release()

	// get/create from the pool
	conns := make([]net.Conn, MaximumCap)
	for i := 0; i < MaximumCap; i++ {
		conn, _ := p.get(context.Background(), network, address)
		conns[i] = conn
	}

	// now put them all back
	for _, conn := range conns {
		p.put(conn)
	}

	if p.len() != MaximumCap {
		t.Errorf("Put error len. Expecting %d, got %d",
			1, p.len())
	}

	p.release() // close pool

}

func TestPool_UsedCapacity(t *testing.T) {
	p, _ := newChannelPool()
	defer p.release()

	if p.len() != InitialCap {
		t.Errorf("InitialCap error. Expecting %d, got %d",
			InitialCap, p.len())
	}
}

func TestPool_Close(t *testing.T) {
	p, _ := newChannelPool()

	// now close it and test all cases we are expecting.
	p.release()

	c := p

	if c.conns != nil {
		t.Errorf("Close error, conns channel should be nil")
	}

	if c.factory != nil {
		t.Errorf("Close error, factory should be nil")
	}

	_, err := p.get(context.Background(), network, address)
	if err == nil {
		t.Errorf("Close error, get conn should return an error")
	}

	if p.len() != 0 {
		t.Errorf("Close error used capacity. Expecting 0, got %d", p.len())
	}
}

func TestPoolConcurrent(t *testing.T) {
	p, _ := newChannelPool()
	pipe := make(chan net.Conn, 0)

	go func() {
		p.release()
	}()

	for i := 0; i < MaximumCap; i++ {
		go func() {
			conn, _ := p.get(context.Background(), network, address)

			pipe <- conn
		}()

		go func() {
			conn := <-pipe
			if conn == nil {
				return
			}
			p.put(conn)
		}()
	}
}

func TestPoolConcurrent2(t *testing.T) {
	//p, _ := NewChannelPool(0, 30, factory)
	p, _ := newChannelPool()

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				conn, _ := p.get(context.Background(), network, address)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				p.close(conn)
				wg.Done()
			}(i)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			conn, _ := p.get(context.Background(), network, address)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			p.close(conn)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func newChannelPool() (*channelPool, error) {
	return NewChannelPool(WithInitialCap(InitialCap), WithMaxCap(MaximumCap), WithFactory(factory), WithClose(closeFac), WithIdleTimeout(time.Second*20))
}

func rpcServer() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", address)
	if e != nil {
		panic(e)
	}
	go http.Serve(l, nil)
}

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
