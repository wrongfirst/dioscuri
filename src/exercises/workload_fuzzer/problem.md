Realistic simulation requires generating synthetic websites and mutating them deterministically across runs.

Given a PRNG, the fuzzer produces reproducible file trees for testing deployment deltas.

## Task

Implement `GenerateSite` and `MutateSite`:

```go
type Site map[string]string

func GenerateSite(rng *rand.Rand, fileCount int) Site
func MutateSite(rng *rand.Rand, base Site, newFiles int) Site
```

- `GenerateSite`:
  - Generate `fileCount` entries: path `fmt.Sprintf("/page_%d.html", i)` with content `fmt.Sprintf("content_%d_%d", i, rng.Intn(1000))`.
- `MutateSite`:
  - Deep-copy `base`.
  - Add `newFiles` entries with path `fmt.Sprintf("/new_%d.html", i)` and content `fmt.Sprintf("new_content_%d", rng.Intn(1000))`.
- Two executions using PRNGs with the same initial seed must produce identical outputs.
