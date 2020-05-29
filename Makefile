default: day

build: bin/h1

bin/h1: cmd/h1/main.go
	go build -o bin/h1 ./cmd/h1

run: build
	./bin/h1 -private -debug

day: bin/h1
	./bin/h1 -private > `date +"%Y-%m-%d-report.md"`
