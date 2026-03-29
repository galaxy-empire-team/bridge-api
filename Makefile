.PHONY: local-run
local-run:
	bash -c 'set -a; . .env; set +a; go run cmd/main.go'

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: gen-transport
gen-transport:
	protoc \
		-I ./api/proto \
  		--go_out=paths=source_relative:./api/gen/go \
  		--go-grpc_out=paths=source_relative:./api/gen/go \
  		api/proto/planet/v1/planet.proto
