In deterministic simulation, invariants are asserted after every step to detect corruption or state leakage immediately.

The 3 core invariants are:
1. **Blob Integrity**: Every blob's content must match its SHA-256 key.
2. **No Dangling References**: Every hash in every manifest must exist in blob storage.
3. **Production Validity**: `ActiveDeployID` must resolve to an existing manifest.

## Task

Implement the invariant verifiers:

```go
var (
    ErrIntegrityViolation = errors.New("blob content does not match hash")
    ErrDanglingReference  = errors.New("manifest references non-existent blob")
    ErrProductionInvalid  = errors.New("production pointer invalid")
)

type SystemState struct {
    ActiveDeployID string
    Manifests      map[string]map[string]string
    Blobs          map[string][]byte
}

func VerifyIntegrity(blobs map[string][]byte) error
func VerifyNoDangling(manifests map[string]map[string]string, blobs map[string][]byte) error
func VerifyProduction(state *SystemState) error
```

- `VerifyIntegrity`: Compute SHA-256 for each blob. Return `ErrIntegrityViolation` if any key does not match its digest.
- `VerifyNoDangling`: Verify every hash in every manifest exists in `blobs`. Return `ErrDanglingReference` if missing.
- `VerifyProduction`: Return `ErrProductionInvalid` if `ActiveDeployID` is empty or missing from `Manifests`.
