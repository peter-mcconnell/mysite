.PHONY: help
help: ## display help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: build ## run the docker container
	@docker run \
		--rm \
		-p 8080:8080 \
		-ti pemcconnell/mysite:latest

.PHONY: build
build: ## build the docker container
	@docker build \
		-t pemcconnell/mysite:latest \
		.
