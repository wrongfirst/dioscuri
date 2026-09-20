Deterministic simulation requires synthesizing workloads with realistic change sets (file creations, modifications, and deletions) while guaranteeing 100% reproducibility from a PRNG seed.

`GenerateSite` and `MutateSite` construct synthetic websites and apply deterministic delta mutations.

## Execution Flow

1. **`GenerateSite(rng, fileCount)`**:
   - Create a `Site` map with `fileCount` entries.
   - Set path to `fmt.Sprintf("/page_%d.html", i)` and content to `fmt.Sprintf("content_%d_%d", i, rng.Intn(1000))`.
2. **`MutateSite(rng, base, toAdd, toModify, toDelete)`**:
   - Clone `base` into a new `Site` map.
   - Sort existing keys alphabetically to guarantee deterministic mutation ordering regardless of Go's map iteration randomization.
   - **Delete**: Remove up to `toDelete` keys from the clone.
   - **Modify**: Overwrite content of up to `toModify` remaining keys with `fmt.Sprintf("modified_%d", rng.Intn(1000))`.
   - **Add**: Insert `toAdd` new entries at `fmt.Sprintf("/new_%d.html", i)` with content `fmt.Sprintf("new_content_%d", rng.Intn(1000))`.
   - Return the mutated `Site`.

## Task

Implement `GenerateSite` and `MutateSite`:

```go
type Site map[string]string

func GenerateSite(rng *rand.Rand, fileCount int) Site
func MutateSite(rng *rand.Rand, base Site, toAdd, toModify, toDelete int) Site
```
