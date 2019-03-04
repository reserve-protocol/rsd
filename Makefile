export REPO_DIR = $(shell pwd)

abi/bindings: contracts/*.sol abi/generate.go
	npx sol-compiler
	go run abi/generate.go
	@echo "placeholder output file for 'make abi/bindings'" > abi/bindings

test: abi/bindings
	go test ./tests

fuzz:
	make -C ./tests fuzz

coverage: abi/bindings
	COVERAGE_ENABLED=1 go test -v ./tests
	open tests/coverage/index.html

fmt:
	npx solium -d contracts/ --fix
	npx solium -d tests/echidna/ --fix

poke: abi/bindings cmd/poke/poke.go
	go install ./cmd/poke

run-dev-container:
	docker run \
		--rm \
		-it \
		--mount type=bind,source="$(REPO_DIR)",target=/reserve-dollar \
		reserveprotocol/env

run-geth:
	docker run -it --rm -p 8545:8501 0xorg/devnet
