abi/bindings: contracts/*.sol abi/generate.go
	npx sol-compiler
	go run abi/generate.go
	touch abi/bindings

test: abi/bindings
	go test ./tests

export REPO_DIR = $(shell pwd)
coverage: abi/bindings
	go test -v -cover ./tests
	open tests/coverage/index.html

fmt:
	npx solium -d contracts/ --fix
