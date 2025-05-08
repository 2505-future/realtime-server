export GOOS=linux
export GOARCH=amd64

LAMBDA_DIRS := $(shell ls lambda)

.PHONY: clean
clean:
	@rm -rf ./bin/*

.PHONY: build
build:
	@for dir in $(LAMBDA_DIRS); do \
		echo "Building $$dir..."; \
		mkdir -p bin/$$dir; \
		go build -o bin/$$dir/bootstrap lambda/$$dir/main.go; \
	done

.PHONY: zip
zip:
	@for dir in $(LAMBDA_DIRS); do \
		echo "Zipping $$dir..."; \
		(cd bin/$$dir && zip ../$$dir.zip bootstrap); \
	done

.PHONY: deploy
deploy: clean build zip
	@echo "Ready to deploy with Terraform or CLI"
