SHELL := /bin/bash
NAME = k8s-inspect

#VERSION?=$(shell git describe --tags --always)
VERSION?=latest

all: clean build

.PHONY: build tinybuild run clean
default: build

build: 
	docker build -t $(NAME):$(VERSION) .

tinybuild: 
	docker build -f Dockerfile.multistage -t $(NAME):$(VERSION) .

run: 
	docker run -p 8080:8080 --name=$(NAME) -d $(NAME):$(VERSION)

clean:
	-docker rm -f $(NAME)
	-docker rmi $(NAME):$(VERSION)
