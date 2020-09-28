SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

doc:
	godoc -http=:6060

shor:
	go run cmd/shor/main.go --N=21 --t=4 --shot=10