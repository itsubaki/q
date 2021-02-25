SHELL := /bin/bash

test:
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

bench:
	cd pkg/math/vector; go test --bench . --benchmem
	cd pkg/math/matrix; go test --bench . --benchmem

doc:
	godoc -http=:6060

shor:
	go run cmd/shor/main.go --N=21 --t=5
