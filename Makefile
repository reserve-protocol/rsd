bindings:
	go run abi/generate.go

test: bindings
	go test ./tests

fmt:
	npx solium -d contracts/ --fix
