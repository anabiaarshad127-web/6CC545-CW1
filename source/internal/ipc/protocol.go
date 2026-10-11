package ipc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"example.edu/6cc545/cw1/internal/config"
)

const MaxMessage = 4096
const Version = 1

var ErrUnimplemented = errors.New("TODO Week 4: implement IPC")

type Message struct {
	Version   int    `json:"version"`
	RequestID string `json:"request_id"`
	NodeID    string `json:"node_id"`
	Text      string `json:"text"`
}

type Ack struct {
	Version   int    `json:"version"`
	RequestID string `json:"request_id"`
	Receiver  string `json:"receiver"`
	Accepted  bool   `json:"accepted"`
}

func (m Message) Validate() error {
	if m.Version != Version {
		return fmt.Errorf("unsupported message version: %d", m.Version)
	}
	if !config.ValidID(m.RequestID) {
		return errors.New("invalid request ID")
	}
	if !config.ValidID(m.NodeID) {
		return errors.New("invalid node ID")
	}
	if m.Text == "" {
		return errors.New("message text cannot be empty")
	}
	if !utf8.ValidString(m.Text) {
		return errors.New("message text must be valid UTF-8")
	}
	return nil
}

func Encode(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 || len(b) > MaxMessage {
		return nil, errors.New("message exceeds limit")
	}
	return b, nil
}

func Decode(data []byte, v any) error {
	if len(data) == 0 {
		return errors.New("empty message")
	}
	if len(data) > MaxMessage {
		return errors.New("message exceeds limit")
	}
	if !utf8.Valid(data) {
		return errors.New("message is not valid UTF-8")
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("message must contain one JSON value")
		}
		return fmt.Errorf("invalid trailing JSON: %w", err)
	}

	if m, ok := v.(*Message); ok {
		return m.Validate()
	}

	return nil
}

func WriteFrame(w io.Writer, data []byte) error {
	if len(data) == 0 || len(data) > MaxMessage {
		return errors.New("invalid frame length")
	}

	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(data)))

	if err := writeAll(w, header[:]); err != nil {
		return fmt.Errorf("write frame header: %w", err)
	}
	if err := writeAll(w, data); err != nil {
		return fmt.Errorf("write frame body: %w", err)
	}
	return nil
}

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		if n <= 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}

func ReadFrame(r io.Reader) ([]byte, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("read frame header: %w", err)
	}

	size := binary.BigEndian.Uint32(header[:])
	if size == 0 || size > MaxMessage {
		return nil, fmt.Errorf("invalid frame length: %d", size)
	}

	body := make([]byte, int(size))
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, fmt.Errorf("read frame body: %w", err)
	}

	return body, nil
}
