# Agent Instructions: go-swe-api

**CRITICAL RULE:** Always run this application and its tests inside the provided Docker container environment. Do NOT attempt to use `go run` or compile binaries locally on the host machine. The project depends on the Swiss Ephemeris C library (`libswe`) and ephemeris data files that are only available inside the containerized environment.

---

## 1. Tech Stack & Architecture

### Tech Stack
- **Language:** Go (with CGo to link against `libswe`)
- **API Protocol:** gRPC & Protocol Buffers
- **Build & Execution:** Makefile, Podman Compose (Docker compatible)
- **Key Dependencies:**
  - `github.com/mshafiee/swephgo` (Go bindings for Swiss Ephemeris)
  - `github.com/jenujari/planets-lib` (Vedic astrology calculations)
  - `google.golang.org/grpc` (gRPC server)

### Project Structure
The application is a gRPC microservice that calculates astrological planetary positions and Shadbala values.

- **`proto/`**: Protocol Buffer definitions (`swe.proto`) and generated Go code.
- **`lib/`**: Core logic for interfacing with the Swiss Ephemeris C library, mapping to domain models (`PlanetCord`, `PlanetBalas`), and calculating Shadbala metrics.
- **`grpc/`**: gRPC server implementation (`server.go`, `root.go`) mapping RPC requests to `lib/` package functions.
- **`client/`**: High-level Go client SDK (`EphServiceClient`) for easy consumption of the gRPC API.

---

## 2. Makefile Commands

The `Makefile` at the root of the project contains all necessary commands to build, test, and interact with the application.

### App Execution & Interaction
- `make sweapi`: Builds the production-ready Docker image and starts the gRPC server in detached mode using `podman compose`.
- `make build-sweapi`: Builds the Docker image for the gRPC API server.
- `make grpc-ui`: Spins up a `grpcui` container on `localhost:8080` that connects to the running gRPC server on `localhost:5678`. Provides a web UI to test and interact with RPC methods.
- `make proto-gen`: Generates Go proto stubs using the `buf` tool inside a container.
- `make down`: Stops and removes all containers and volumes defined in `compose.yaml`.

### Testing & Validation
Because of the C-library dependency, tests must be run via the Makefile:

```bash
# Run all tests in the containerized environment
make sweapi-test

# Run a specific test by name (accepts Go test regex)
make sweapi-test TEST=Test_GetAllPlanetsBalas
```

**Host Validation Escape Hatch:**
If you only need to check that Go code compiles and passes static analysis (no linking or execution required), you can run:
```bash
go vet ./...
```
This skips the C linker step and is safe to run on the host.

---

## 3. Current Application State

- **Implemented RPCs:** The gRPC server currently exposes endpoints for `Ping`, `GetPos`, `FindConjunction`, `Tithy`, and `GetBalas`.
- **Core Engine (`lib/`):** Functions like `GetPlanetCalculation`, `GetAllPlanetsBalas`, and `FindConjunctionRange` are fully implemented and unit-tested.
- **Recent Refactoring:** `GetAllPlanetsBalas` was recently refactored for better readability and safety (now returning a concrete `map[string]PlanetBalas`).
- **New Feature (`GetBalas`):** The `GetBalas` RPC and client SDK method were recently added. It enables fetching detailed planetary coordinates, speeds, and 5 shadbala metrics (Uday, Uchcha, Vakra, Kshetra, Navamsha) for all 9 vedic planets over gRPC.
- **Dependency Management:** Unused modules like the `goforj/godump` package are managed safely using a `tools.go` build constraint to prevent pruning by `go mod tidy`.
