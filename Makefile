SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /cmd/ ) -v -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o coverage.html

lint:
	golangci-lint run

vet:
	go vet ./...

bench:
	go test -bench . ./... --benchmem

doc:
	godoc -http=:6060

shor:
	go run cmd/shor/main.go --N 21 -t 5

grover:
	go run cmd/grover/main.go -top 8

counting:
	go run cmd/counting/main.go -t 7 -top 8
