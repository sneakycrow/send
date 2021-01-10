.PHONY: publish \
				build

GIT_COMMIT    := $(shell git rev-parse --short HEAD)
VERSION				:= ${GIT_COMMIT}
IMAGE					:= sneakycrow/send
REGISTRY			:= ghcr.io

publish:
	docker push ${REGISTRY}/${IMAGE}:latest
	docker push ${REGISTRY}/${IMAGE}:${VERSION}


build:
	docker build --build-arg VERSION=${VERSION} -t ${IMAGE}:${VERSION} .
	docker tag ${IMAGE}:${VERSION} ${REGISTRY}/${IMAGE}:${VERSION}
	docker tag ${IMAGE}:${VERSION} ${REGISTRY}/${IMAGE}:latest

start: 
	docker run --name send -p 8080:3000 -v "$(pwd)/tmp":/tmp -d sneakycrow/send:latest