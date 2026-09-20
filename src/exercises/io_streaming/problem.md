Reading entire files into memory with `os.ReadFile` or `io.ReadAll` exhausts memory when dealing with large payloads.

In Go, streaming between an `io.Reader` and an `io.Writer` avoids loading the full byte stream into RAM. The standard `io.Copy(dst, src)` function transfers data in fixed-size chunk buffers until `io.EOF`.

## Task

Implement `CopyAndCount`:

```go
func CopyAndCount(dst io.Writer, src io.Reader) (int64, error)
```

- Stream all bytes from `src` into `dst`.
- Return the total number of bytes transferred.
- Return any error encountered during the operation.
