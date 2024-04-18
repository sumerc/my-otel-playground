BINARY_NAME=my-otel-collector

SOURCE_FILE=main.go

build:
	go build -o $(BINARY_NAME) $(SOURCE_FILE)
run: build
	./$(BINARY_NAME)
clean:
	rm -f $(BINARY_NAME)

.PHONY: build clean
