In simulation testing, an invariant oracle validates system state after every step to detect data corruption, dangling pointers, or broken invariants immediately.

`Oracle` inspects the filesystem and blob store directly to assert the three core invariants:
1. **Blob Integrity**: Every blob's content matches its SHA-256 filename.
2. **Production Safety**: `/production.json` points to an existing, valid manifest.
3. **No Dangling References**: Every hash referenced in every manifest exists in the blob store.

## Execution Flow

1. **`AssertBlobIntegrity()`**:
   - List files in `store.blobDir`.
   - Stream each blob via `store.OpenBlob(hash)` through `crypto/sha256`.
   - If the computed digest does not match the filename, return `ErrIntegrityViolation`.
2. **`AssertProductionSafety()`**:
   - Open `/production.json` and decode `{"deploy_id": "..."}`.
   - Return `ErrProductionInvalid` if missing, empty, or if `/manifests/<deploy_id>.json` does not exist on `fs`.
3. **`AssertNoDanglingReferences()`**:
   - List all `*.json` files in `manifestDir`.
   - Decode each manifest (`map[string]string`) and check `store.HasBlob(hash)` for every value.
   - Return `ErrDanglingReference` if any referenced blob is absent.
4. **`AssertAll()`**:
   - Execute the three assertions. Return the first error encountered, or `nil` if all invariants hold.

## Task

Implement `Oracle`:

```go
var (
    ErrIntegrityViolation = errors.New("blob content does not match hash")
    ErrDanglingReference  = errors.New("manifest references non-existent blob")
    ErrProductionInvalid  = errors.New("production pointer invalid")
)

type Oracle struct {
    fs          FS
    store       *Store
    manifestDir string
}

func NewOracle(fs FS, store *Store, manifestDir string) *Oracle
func (o *Oracle) AssertBlobIntegrity() error
func (o *Oracle) AssertProductionSafety() error
func (o *Oracle) AssertNoDanglingReferences() error
func (o *Oracle) AssertAll() error
```
