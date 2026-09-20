In Content-Addressable Storage (CAS), data is identified by the cryptographic hash of its content rather than its filename or path.

In Go, `crypto/sha256.New()` returns a hash instance implementing `io.Writer`. Bytes can be streamed directly from an `io.Reader` into the hasher using `io.Copy` without buffering the payload in memory.

## Task

Implement `HashReader`:

```go
func HashReader(r io.Reader) (string, error)
```

- Stream all bytes from `r` into a SHA-256 hasher.
- Return the 64-character lowercase hexadecimal digest string.
- If reading fails, return an empty string and the error.
