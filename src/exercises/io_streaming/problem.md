Reading entire files into memory with `os.ReadFile` exhausts RAM when processing large files. To achieve constant $O(1)$ memory usage, data must be transferred in fixed-size chunks from an `io.Reader` into an `io.Writer`.

## Execution Flow

1. **Validation**: Check that `bufSize > 0`. If not, return `0, ErrInvalidBufferSize`.
2. **Buffer Allocation**: Allocate a single reusable chunk buffer: `make([]byte, bufSize)`.
3. **Transfer Loop**:
   - Call `src.Read(buf)` to fill the buffer.
   - If `nr > 0` bytes were read, call `dst.Write(buf[:nr])`.
   - Add bytes written to the running total. If `nw != nr`, return the total and `io.ErrShortWrite`.
   - Check the error from `src.Read`:
     - If `io.EOF`, all bytes have been transferred. Break cleanly.
     - For any other error, return immediately.

## Task

Implement `StreamCopy`:

```go
func StreamCopy(dst io.Writer, src io.Reader, bufSize int) (int64, error)
```

- Stream all bytes from `src` to `dst` using chunks of `bufSize`.
- Return the total byte count written and any error encountered (`nil` on success).
