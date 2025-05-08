export GOOS=linux
export GOARCH=amd64

LAMBDA_DIRS := $(shell ls lambda)

DATE := $(shell TZ=Asia/Tokyo date +%Y%m%d)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo $(DATE))

ECR_REPO := 794038226787.dkr.ecr.ap-northeast-1.amazonaws.com

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

push-connect:
	docker build -f ./docker/Dockerfile.connect -t 58hack-connect .
	docker tag 58hack-connect:latest $(ECR_REPO)/58hack-connect:$(COMMIT)
	docker push $(ECR_REPO)/58hack-connect:$(COMMIT)

push-disconnect:
	docker build -f ./docker/Dockerfile.disconnect -t 58hack-disconnect .
	docker tag 58hack-disconnect:latest $(ECR_REPO)/58hack-disconnect:$(COMMIT)
	docker push $(ECR_REPO)/58hack-disconnect:$(COMMIT)

push-send:
	docker build -f ./docker/Dockerfile.send -t 58hack-send .
	docker tag 58hack-send:latest $(ECR_REPO)/58hack-send:$(COMMIT)
	docker push $(ECR_REPO)/58hack-send:$(COMMIT)
