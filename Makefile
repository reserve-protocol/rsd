export REPO_DIR = $(shell pwd)

abi/bindings: contracts/*.sol abi/generate.go
	npx sol-compiler
	go run abi/generate.go
	@echo "placeholder output file for 'make abi/bindings'" > abi/bindings

test: abi/bindings
	go test ./tests

coverage: abi/bindings
	go test -v -cover ./tests
	open tests/coverage/index.html

fmt:
	npx solium -d contracts/ --fix

res: abi/bindings cmd/res/res.go
	go install ./cmd/res

run-dev-container:
	docker run \
		--rm \
		-it \
		--mount type=bind,source="$(REPO_DIR)",target=/reserve-dollar \
		reserveprotocol/env
