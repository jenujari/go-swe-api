PHONY: all

TEST ?= .

down:
	podman compose -f compose.yaml down --remove-orphans --volumes
	@echo "Stopped and removed all containers"

build-swe-base:
	podman build -f swe-builder.Dockerfile -t docker.io/jhon5456/sweph-build-base:v2.10.03 .
	podman push docker.io/jhon5456/sweph-build-base:v2.10.03
	@echo "Built and pushed sweph-build-base image"

# Generate proto code using podman and buf (cleaner and more reliable)
proto-gen:
	podman run --rm -v .:/workspace -w /workspace docker.io/bufbuild/buf generate proto
	@echo "Generated proto code"

build-sweapi-test: proto-gen
	podman compose -f compose.yaml build test_sweapi
	@echo "Built test_sweapi image"


# make sweapi-test TEST=PosHandler
sweapi-test:
	@echo "Running test $(TEST)"
	podman compose -f compose.yaml run --rm test_sweapi ./... -v -run $(TEST) -coverpkg=./... -coverprofile=coverage.out
	podman compose -f compose.yaml run --entrypoint sh --rm test_sweapi -c "grep -v '\\.pb\\.go' coverage.out > cov.tmp && mv cov.tmp coverage.out && go tool cover -func=coverage.out"
	@echo "Test $(TEST) completed"


build-sweapi: proto-gen
	podman compose -f compose.yaml build sweapi
	@echo "Built sweapi image"

sweapi: build-sweapi
	podman compose -f compose.yaml up -d sweapi

grpc-ui:
	podman run --rm --network=host -p 8080:8080 docker.io/fullstorydev/grpcui -plaintext localhost:5678

build-and-push-container:
	./scripts/build_and_push.sh

jenkins-deploy:
	./scripts/trigger_deploy.sh