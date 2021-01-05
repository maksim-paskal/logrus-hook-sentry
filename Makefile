test:
	./scripts/validate-license.sh
	go fmt
	go mod tidy
	go test -race
	golangci-lint run --allow-parallel-runners -v --enable-all --disable exhaustivestruct,testpackage --fix
testScript:
	go run ./test-script
upgrade:
	go get -v -u all
	go mod tidy