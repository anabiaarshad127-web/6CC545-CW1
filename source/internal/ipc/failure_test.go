package ipc

import (
	"context"
	"encoding/binary"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestWeek04_ServerRejectsOversizeTCP(t *testing.T) {
	ln, e := net.Listen("tcp", "127.0.0.1:0")
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer ln.Close()
	var calls atomic.Int32
	done := make(chan error, 1)
	go func() { done <- ServeTCP(ctx, ln, "node-b", 100*time.Millisecond, func(Message) { calls.Add(1) }) }()
	c, e := net.DialTimeout("tcp", ln.Addr().String(), 200*time.Millisecond)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(300 * time.Millisecond))
	var h [4]byte
	binary.BigEndian.PutUint32(h[:], MaxMessage+1)
	_, _ = c.Write(h[:])
	var b [1]byte
	n, e := c.Read(b[:])
	if e == nil || n != 0 {
		t.Error("oversize TCP frame was not rejected")
	}
	if calls.Load() != 0 {
		t.Error("rejected frame mutated inbox")
	}
	cancel()
	select {
	case e := <-done:
		if e != nil {
			t.Errorf("serve failed: %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("shutdown hung")
	}
}
func TestWeek04_UDPDropAndDeadline(t *testing.T) {
	c, e := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer c.Close()
	done := make(chan error, 1)
	var calls atomic.Int32
	go func() { done <- ServeUDP(ctx, c, "node-b", 100*time.Millisecond, func(Message) { calls.Add(1) }) }()
	remote, e := net.DialTimeout("udp", c.LocalAddr().String(), time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer remote.Close()
	_ = remote.SetDeadline(time.Now().Add(150 * time.Millisecond))
	_, _ = remote.Write([]byte(`{malformed}`))
	var b [128]byte
	n, e := remote.Read(b[:])
	if e == nil || n != 0 {
		t.Error("malformed UDP request should be dropped")
	}
	if calls.Load() != 0 {
		t.Error("rejected datagram changed inbox")
	}
	cancel()
	select {
	case e := <-done:
		if e != nil {
			t.Errorf("serve: %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("UDP shutdown hung")
	}
}

func TestWeek04_UDPRejectsOversizeDatagram(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer conn.Close()

	var calls atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- ServeUDP(ctx, conn, "node-b", 100*time.Millisecond, func(Message) {
			calls.Add(1)
		})
	}()

	remote, err := net.DialTimeout("udp", conn.LocalAddr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer remote.Close()
	if err := remote.SetDeadline(time.Now().Add(250 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	if _, err := remote.Write(make([]byte, MaxMessage+1)); err != nil {
		t.Fatal(err)
	}

	var response [128]byte
	n, err := remote.Read(response[:])
	if err == nil || n != 0 {
		t.Errorf("oversize UDP datagram should be dropped; n=%d err=%v", n, err)
	}
	if calls.Load() != 0 {
		t.Error("oversize UDP datagram changed the inbox")
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("UDP server shutdown: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("UDP server shutdown hung")
	}
}

func TestWeek04_TCPIncompleteFrameTimesOut(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer ln.Close()

	var calls atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- ServeTCP(ctx, ln, "node-b", 100*time.Millisecond, func(Message) {
			calls.Add(1)
		})
	}()

	client, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if err := client.SetDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}

	var header [4]byte
	binary.BigEndian.PutUint32(header[:], 10)
	if _, err := client.Write(header[:]); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Write([]byte{'{'}); err != nil {
		t.Fatal(err)
	}

	var response [1]byte
	n, err := client.Read(response[:])
	if err == nil || n != 0 {
		t.Errorf("incomplete TCP frame should time out without an acknowledgement; n=%d err=%v", n, err)
	}
	if calls.Load() != 0 {
		t.Error("incomplete TCP frame changed the inbox")
	}

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("TCP server shutdown: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("TCP server shutdown hung")
	}
}
