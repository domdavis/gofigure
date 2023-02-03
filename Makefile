all: clean build test vet lint

ci: all report

clean:
	go clean
	rm -f *.out *.prof *.test
	go mod tidy

build:
	go build

vet:
	go vet ./...

lint:
	golangci-lint run
	golint ./...

test:
	go test -v -covermode=count ./...

doc:
	godoc --http=:6060

profile:
	go test --cpuprofile cpu.prof --memprofile mem.prof -bench ./...

report:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
