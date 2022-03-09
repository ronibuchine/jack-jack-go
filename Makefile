
MAIN=./cmd/compiler/compiler.go

default:
	go build -o jjg $(MAIN)

test:
	go test ./...
