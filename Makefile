GO111MODULE := on
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

all: alertlogger

.PHONY: alertlogger
alertlogger:
	go build -o bin/alertlogger
	strip bin/alertlogger

.PHONY: dockerimages
dockerimages:
	docker build -t mwennrich/alertlogger:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush:
	docker push mwennrich/alertlogger:${DOCKER_TAG}

.PHONY: clean
clean:
	rm -f bin/*

