bindings:
	go run abi/generate.go

test: bindings
	go test ./tests

export REPO_DIR = $(shell pwd)
coverage: bindings
	go test -v -cover ./tests
	open tests/coverage/index.html

fmt:
	npx solium -d contracts/ --fix
