In Content-Addressable Storage (CAS), objects are identified by the cryptographic hash of their content rather than arbitrary filenames. This guarantees tamper-proof integrity and automatic deduplication.

Go's `crypto/sha256.New()` creates an instance of `hash.Hash`, which implements `io.Writer`. This allows streaming bytes directly from an `io.Reader` into the hasher with constant $O(1)$ memory.

## Execution Flow

1. **Initialize Hasher**: Call `sha256.New()` to create the SHA-256 writer.
2. **Stream Content**: Stream bytes from `r` into the hasher using `io.Copy(hasher, r)`.
3. **Handle Errors**: If `io.Copy` encounters an error, return an empty string and the error.
4. **Finalize Digest**: Call `hasher.Sum(nil)` to retrieve the 32-byte binary digest.
5. **Hex Encode**: Format the digest into a 64-character lowercase hexadecimal string using `hex.EncodeToString`.

## Task

Implement `HashReader`:

```go
func HashReader(r io.Reader) (string, error)
```

- Stream all bytes from `r` through SHA-256.
- Return the hexadecimal digest string and any error encountered.

### Test Vectors
- Empty stream `""` $\to$ `"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`
- `"hello world"` $\to$ `"b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"`
