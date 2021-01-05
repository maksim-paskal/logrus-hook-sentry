test:
	./scripts/validate-license.sh
	go fmt
	go mod tidy
	go test -race
	golangci-lint run --allow-parallel-runners -v --enable-all --disable exhaustivestruct --fix