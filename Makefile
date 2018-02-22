SHELL := /bin/bash
NAME = goinspect

#VERSION?=$(shell git describe --tags --always)
VERSION?=latest

all: docker

.PHONY: docker run
default: build

docker: 
	docker build -t $(NAME):$(VERSION) .

run: 
	docker run -p 8080:8080 --name=$(NAME) -d $(NAME):$(VERSION)

clean:
	-docker rm -f $(NAME)
	-docker rmi $(NAME):$(VERSION)
