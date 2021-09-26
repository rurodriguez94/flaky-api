run:
	go run main.go -pages=10 -retries=10 -stopOnFail=false

build-run:
	go build
	./flaky-api -pages=10 -retries=10 -stopOnFail=false

test:
	go test -v ./...

remove-images:
	rm -r images