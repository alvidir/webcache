VERSION=0.2.0
PROJECT=webcache
REPO=alvidir
REMOTE=docker.io

release: build push

build:
	podman build -t ${REPO}/${PROJECT}:${VERSION} -f ./docker/webcache/dockerfile .

push:
	podman tag localhost/${REPO}/${PROJECT}:${VERSION} ${REMOTE}/${REPO}/${PROJECT}:${VERSION}
	podman push ${REMOTE}/${REPO}/${PROJECT}:${VERSION}

deploy:
	podman-compose -f docker-compose.yaml up --remove-orphans -d

undeploy:
	podman-compose -f docker-compose.yaml down

run:
	go run cmd/webcache/main.go

test:
	go test -v -race ./...