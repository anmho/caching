
.PHONY: build
build:
	go build -o ./bin/api/api main.go

.PHONY: run
run: build
