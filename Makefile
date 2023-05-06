.PHONY: lint


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: audit
audit: module lint test

.PHONY: module
module:
	@printf "\e[1mTidying and verifying module dependencies...\e[0m\n"
	@go mod tidy
	@go mod verify

.PHONY: lint
lint:
	@golangci-lint --version > /dev/null 2>&1;\
    	if [ "$$?" != 0 ]; then\
    		printf "\e[31mgolangci-lint is not installed. Run go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0\e[0m\n";\
    		exit 1;\
    	fi
	@printf "\e[1mRunning \e[36mgolangci-lint\e[0m\n"
	@golangci-lint run

.PHONY: test
test:
	@printf "\e[1mExecuting all test files...\e[0m\n"
	@go test -race -cover ./...