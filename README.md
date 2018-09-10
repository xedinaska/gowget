# Gowget

Gowget is a simple app that's very similar to wget. Use it to download one or more files by provided URLs.

## Makefile

You can build binary using Makefile. Just type `make build` and binary will appear in `./bin` folder.

You can format & lint your code using make lint command. It includes following apps: `gofmt, golint, govet, gocyclo`.

## Usage

Build binary file using `make build` command.

Provide files that should be downloaded as args: `./bin/gowget "http://first-url.com" "https://second-url.com"`.
All files will be downloaded to `./etc/downloads` folder.

In case of any error you'll be prompted.

## TODO

- Move `http.*` calls out of  the downloader package (inject dependencies for easy mocking in tests)
- Same for `io.*` calls
- Cover with unit tests & add appropriate task to Makefile