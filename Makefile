
.PHONY: build
build:
	@go build -o ./bin/api main.go

.PHONY: run
run: build
	@./bin/api

.PHONY: watch
watch:
	@air --build.cmd "make build" --build.bin "./bin/api"

.PHONY: test
test:
	@go test ./...