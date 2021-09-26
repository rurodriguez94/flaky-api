# Flaky API

---

## Description
This project works like a script, when run find a list of houses and downloads their images and save them in the folder named images that is created automatically.

---

## Run application:

It can be run in two ways, directly with the main file or executing its compiled binary,
in both cases you can use the makefile to run the commands.

This application can use three non-mandatory arguments:
* pages: used to define the number of pages to fetch (default value is 10)
* retries: used to define the number of http retries (default value is 10)
* stopOnFail: stop the program execution when the max retries of any request is reached (default value is false)

### Commands to run the app:

* Using commands directly:

        Using default arguments:
        go run main.go
  
        Using arguments:
        go run main.go -pages=10 -retries=10 -stopOnFail=false
        
        Using the binary builded
        go build && ./flaky-api -pages=10 -retries=10 -stopOnFail=false

* Using Makefile:
       
        make run (execute: go run main.go -pages=10 -retries=10 -stopOnFail=false)
        make build-run (build the binary and execute it)
---

## Library used:

In this project 3 libraries was used:

* https://github.com/go-resty/resty
      
      Provides a simple http client to make request (making it easier) and have implemented
      a exponential backoff alghoritm for retries (implies exponential backoff with jitter) by default
      based in: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

* https://github.com/stretchr/testify

      Used for unit tests (provide asserts and require functions)

* https://github.com/sirupsen/logrus

      Used for logging using structured logger.

---

## Project layout:

For the project layout, I used as a reference:

* https://github.com/golang-standards/project-layout
* https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html

---

## Other commands:

    make test (run all unit tests)
    make remove-images (remove the images folder)

---