PHONY: all

TEST ?= .

build-swe-base:
	podman build -f swe-builder.Dockerfile -t swe-builder:latest . 
	echo "Built swe-builder image"

build-sweapi-test:
	podman compose -f compose.dev.yaml build test_sweapi 
	echo "Built test_sweapi image"

# make sweapi-test TEST=PosHandler
sweapi-test:
	echo "Running test $(TEST)"
	podman compose -f compose.dev.yaml  run --rm test_sweapi -run $(TEST) -v 
	echo "Test $(TEST) completed"