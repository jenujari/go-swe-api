# Agent Instructions

## Running Tests

**IMPORTANT:** Do NOT run `go test` directly from the project root or any subdirectory. The project depends on the Swiss Ephemeris C library (`libswe`) and ephemeris data files that are only available inside the containerised test environment.

Always use the Makefile targets to run tests:

```bash
# Run all tests
make sweapi-test

# Run a specific test by name
make sweapi-test TEST=Test_GetAllPlanetsBalas
```

The `TEST` variable accepts a Go test name regex, so you can target individual tests or groups.

### Why?

Running `go test ./...` from the host will fail at the linker stage with:

```
/usr/bin/ld: cannot find -lswe: No such file or directory
```

The `make sweapi-test` target uses `podman compose` to spin up a container that has `libswe` installed, the ephemeris files mounted, and the correct environment configured.

## Tech Stack

- **Language:** Go (with CGo — links against the Swiss Ephemeris C library)
- **Build / Test runner:** Makefile → Podman Compose
- **Proto generation:** `make proto-gen` (uses `buf` inside a container)
- **Key dependencies:**
  - `github.com/mshafiee/swephgo` — Go bindings for Swiss Ephemeris
  - `github.com/jenujari/planets-lib` — Vedic astrology calculations (signs, nakshatras, balas, speed categories)
  - `github.com/stretchr/testify` — Test assertions
  - `google.golang.org/grpc` + `google.golang.org/protobuf` — gRPC server

## Code Validation Without the Container

If you only need to check that Go code compiles and passes static analysis (no linking or execution), you can run:

```bash
go vet ./...
```

This skips the C linker step and will catch syntax errors, type mismatches, and common bugs.
