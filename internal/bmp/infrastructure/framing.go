package infrastructure

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Frame struct {
	Version byte
	Length  uint32
	Payload []byte
}

func ReadFrame(r io.Reader, max uint32) (Frame, error) {
	h := make([]byte, 5)
	if _, e := io.ReadFull(r, h); e != nil {
		return Frame{}, e
	}
	n := binary.BigEndian.Uint32(h[1:])
	if n > max || n < 5 {
		return Frame{}, fmt.Errorf("invalid BMP frame length %d", n)
	}
	p := make([]byte, n-5)
	if _, e := io.ReadFull(r, p); e != nil {
		return Frame{}, e
	}
	return Frame{h[0], n, p}, nil
}
func WriteFrame(w io.Writer, f Frame) error {
	if len(f.Payload)+5 > int(^uint32(0)) {
		return fmt.Errorf("frame too large")
	}
	h := make([]byte, 5)
	h[0] = f.Version
	binary.BigEndian.PutUint32(h[1:], uint32(len(f.Payload)+5))
	if _, e := w.Write(h); e != nil {
		return e
	}
	_, e := w.Write(f.Payload)
	return e
}
