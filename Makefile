all: build

.PHONY: deps
deps:
	dep ensure
	git submodule update --init

.PHONY: build
build: deps
	go build ./...

.PHONY: testenv
testenv:
	( cd testenv ; make testenv )

.PHONY: test
test: vendor
	go test -race -cover ./...

.PHONY: lint
lint: vendor
	golint $(go list ./... | grep -v /vendor/)

.PHONY: clean
clean:
	( cd testenv ; make clean )

.PHONY: check
check:
	go vet ./...
	go test -run xxxx ./...