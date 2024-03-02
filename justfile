build:
	go build -o bin/main main.go

run:
	go run main.go

install:
	go install -ldflags "-s -w -X github.com/gofsd/fsd/cmd.Version=`git tag | tail -n 1` -X github.com/gofsd/fsd/cmd.Build=`git rev-parse HEAD`"

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/main-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all