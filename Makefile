.PHONY: build
build:
	go build -o bin/constanta-test ./cmd/constanta-test/main.go

.PHONY: run
run:
	go run cmd/constanta-test/main.go
