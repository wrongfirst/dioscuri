package main

import (
	"bytes"
	"errors"
	"io"
)

var ErrInvalidBufferSize = errors.New("invalid buffer size")

// StreamCopy streams data from src to dst using a buffer of size bufSize.
// It returns total bytes written and any error encountered.
func StreamCopy(dst io.Writer, src io.Reader, bufSize int) (int64, error) {
	if bufSize <= 0 {
		return 0, ErrInvalidBufferSize
	}

	buf := make([]byte, bufSize)
	var total int64

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[:nr])
			if nw > 0 {
				total += int64(nw)
			}
			if ew != nil {
				return total, ew
			}
			if nr != nw {
				return total, io.ErrShortWrite
			}
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			return total, er
		}
	}

	return total, nil
}
