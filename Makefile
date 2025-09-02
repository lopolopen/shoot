.PHONY: test

test:
	go test ./cmd

testup:
	go test ./cmd -update
