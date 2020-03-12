PACKAGES=$(shell go list ./...)
export GO111MODULE = on

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/irisnet/irishub-sdk-go

test-unit:
	cd test/scripts/ && sh build.sh && sh start.sh
	sleep 5s
	@go test -p 1 $(PACKAGES)
	sh clean.sh