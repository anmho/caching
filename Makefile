
.PHONY: build
build:
	@go build -o ./bin/api ./cmd/api

.PHONY: run
run: build
	@./bin/api

.PHONY: watch
watch:
	@air --build.cmd "make build" --build.bin "./bin/api"

.PHONY: test
test:
	@go test ./...

.PHONY: tf
tf:
	@terraform -chdir=./terraform apply
