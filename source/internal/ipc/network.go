package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

type OnMessage func(Message)

func deadline(ctx context.Context, timeout time.Duration) time.Time {
	d := time.Now().Add(timeout)
	if cd, ok := ctx.Deadline(); ok && cd.Before(d) {
		return cd
	}
	return d
}

func makeAck(id string, m Message) Ack {
	return Ack{
		Version:   Version,
		RequestID: m.RequestID,
		Receiver:  id,
		Accepted:  true,
	}
}

func checkAck(a Ack, m Message) error {
	if a.Version != Version ||
		a.RequestID != m.RequestID ||
		!a.Accepted ||
		a.Receiver == "" {
		return errors.New("invalid or unrelated acknowledgement")
	}
	return nil
}

func ServeTCP(ctx context.Context, ln net.Listener, id string, timeout time.Duration, on OnMessage) error {
	if ctx.Err() != nil {
		return nil
	}
	if timeout <= 0 {
		return errors.New("TCP timeout must be positive")
	}

	stop := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = ln.Close()
		case <-stop:
		}
	}()
	defer close(stop)

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept TCP connection: %w", err)
		}

		handleDone := make(chan struct{})
		go func(c net.Conn) {
			select {
			case <-ctx.Done():
				_ = c.Close()
			case <-handleDone:
			}
		}(conn)

		_ = handleTCP(conn, id, timeout, on)
		close(handleDone)
		_ = conn.Close()

		if ctx.Err() != nil {
			return nil
		}
	}
}

func handleTCP(conn net.Conn, id string, timeout time.Duration, on OnMessage) error {
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}

	data, err := ReadFrame(conn)
	if err != nil {
		return fmt.Errorf("read TCP request: %w", err)
	}

	var m Message
	if err := Decode(data, &m); err != nil {
		return fmt.Errorf("decode TCP request: %w", err)
	}

	if on != nil {
		on(m)
	}

	data, err = Encode(makeAck(id, m))
	if err != nil {
		return fmt.Errorf("encode TCP acknowledgement: %w", err)
	}
	if err := WriteFrame(conn, data); err != nil {
		return fmt.Errorf("write TCP acknowledgement: %w", err)
	}

	return nil
}

func SendTCP(ctx context.Context, address string, m Message, timeout time.Duration) (Ack, error) {
	if err := ctx.Err(); err != nil {
		return Ack{}, err
	}
	if timeout <= 0 {
		return Ack{}, errors.New("TCP timeout must be positive")
	}
	if err := m.Validate(); err != nil {
		return Ack{}, fmt.Errorf("invalid TCP message: %w", err)
	}

	data, err := Encode(m)
	if err != nil {
		return Ack{}, fmt.Errorf("encode TCP request: %w", err)
	}

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return Ack{}, fmt.Errorf("connect TCP: %w", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(deadline(ctx, timeout)); err != nil {
		return Ack{}, err
	}

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	if err := WriteFrame(conn, data); err != nil {
		return Ack{}, fmt.Errorf("send TCP request: %w", err)
	}

	response, err := ReadFrame(conn)
	if err != nil {
		if ctx.Err() != nil {
			return Ack{}, ctx.Err()
		}
		return Ack{}, fmt.Errorf("read TCP acknowledgement: %w", err)
	}

	var ack Ack
	if err := Decode(response, &ack); err != nil {
		return Ack{}, fmt.Errorf("decode TCP acknowledgement: %w", err)
	}
	if err := checkAck(ack, m); err != nil {
		return Ack{}, err
	}

	return ack, nil
}

func ServeUDP(ctx context.Context, conn *net.UDPConn, id string, timeout time.Duration, on OnMessage) error {
	if ctx.Err() != nil {
		return nil
	}
	if timeout <= 0 {
		return errors.New("UDP timeout must be positive")
	}

	stop := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-stop:
		}
	}()
	defer close(stop)

	buffer := make([]byte, MaxMessage+1)

	for {
		n, remote, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read UDP datagram: %w", err)
		}

		if n == 0 || n > MaxMessage {
			continue
		}

		var m Message
		if err := Decode(buffer[:n], &m); err != nil {
			continue
		}

		if on != nil {
			on(m)
		}

		data, err := Encode(makeAck(id, m))
		if err != nil {
			continue
		}

		if err := conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
			continue
		}
		if _, err := conn.WriteToUDP(data, remote); err != nil && ctx.Err() != nil {
			return nil
		}
	}
}

func SendUDP(ctx context.Context, address string, m Message, timeout time.Duration) (Ack, error) {
	if err := ctx.Err(); err != nil {
		return Ack{}, err
	}
	if timeout <= 0 {
		return Ack{}, errors.New("UDP timeout must be positive")
	}
	if err := m.Validate(); err != nil {
		return Ack{}, fmt.Errorf("invalid UDP message: %w", err)
	}

	data, err := Encode(m)
	if err != nil {
		return Ack{}, fmt.Errorf("encode UDP request: %w", err)
	}

	remote, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return Ack{}, fmt.Errorf("resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, remote)
	if err != nil {
		return Ack{}, fmt.Errorf("connect UDP: %w", err)
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	if err := conn.SetDeadline(deadline(ctx, timeout)); err != nil {
		return Ack{}, err
	}
	if _, err := conn.Write(data); err != nil {
		if ctx.Err() != nil {
			return Ack{}, ctx.Err()
		}
		return Ack{}, fmt.Errorf("send UDP request: %w", err)
	}

	buffer := make([]byte, MaxMessage+1)
	n, err := conn.Read(buffer)
	if err != nil {
		if ctx.Err() != nil {
			return Ack{}, ctx.Err()
		}
		return Ack{}, fmt.Errorf("read UDP acknowledgement: %w", err)
	}
	if n == 0 || n > MaxMessage {
		return Ack{}, errors.New("invalid UDP acknowledgement size")
	}

	var ack Ack
	if err := Decode(buffer[:n], &ack); err != nil {
		return Ack{}, fmt.Errorf("decode UDP acknowledgement: %w", err)
	}
	if err := checkAck(ack, m); err != nil {
		return Ack{}, err
	}

	return ack, nil
}
