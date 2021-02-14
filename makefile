test:
	go test ./...

test-cover:
ifeq ($(wildcard coverage/),)
	mkdir coverage
endif
	go test ./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	go tool cover -func=coverage/coverage.out