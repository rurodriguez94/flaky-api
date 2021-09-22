run:
	go run main.go -pages=10 -retries=10

build-run:
	go build
	./flaky-api -pages=10 -retries=10

remove-images:
	rm -r images