TARGET=meshbird

all: clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -v

build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(TARGET)
	GOOS=windows GOARCH=amd64 go build -o bin/windows/$(TARGET).exe
	GOOS=darwin GOARCH=amd64 go build -o bin/mac/$(TARGET)

fmt:
	go fmt *.go
