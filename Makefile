all: build

# Build manager binary
build: 
	go build -o bin/getxpath main.go

