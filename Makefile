export GOLANG_PROTOBUF_REGISTRATION_CONFLICT := warn

.PHONY: clean
clean:
	rm -rf gen
	rm -rf dist

.PHONY: update
update:
	docker pull ghcr.io/goccy/bigquery-emulator:0.6.6 --platform linux/amd64
	go mod tidy

.PHONY: test
test: clean update
	ginkgo ./...