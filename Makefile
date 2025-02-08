export GOLANG_PROTOBUF_REGISTRATION_CONFLICT := warn

.PHONY: clean
clean:
	rm -rf gen
	rm -rf dist

.PHONY: update
update:
	go mod tidy

.PHONY: test
test: clean update
	ginkgo ./...