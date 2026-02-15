PHONY: all

TEST ?= .

build-swe-base:
	podman build -f swe-builder.Dockerfile -t swe-builder:latest . 
	echo "Built swe-builder image"

# Generate proto code using podman and buf (cleaner and more reliable)
proto-gen:
	podman run --rm -v .:/workspace -w /workspace docker.io/bufbuild/buf generate proto
	echo "Generated proto code"

build-sweapi-test: build-swe-base proto-gen
	podman compose -f compose.dev.yaml build test_sweapi 
	echo "Built test_sweapi image"



# make sweapi-test TEST=PosHandler
sweapi-test:
	echo "Running test $(TEST)"
	podman compose -f compose.dev.yaml  run --rm test_sweapi -run $(TEST) -v 
	echo "Test $(TEST) completed"