VERSION=1.0.0
PROJECT=webcache
REPO=alvidir
REMOTE=docker.io

release: build push

build:
	podman build -t ${REPO}/${PROJECT}:${VERSION} -f ./container/webcache/containerfile .

push:
	podman tag localhost/${REPO}/${PROJECT}:${VERSION} ${REMOTE}/${REPO}/${PROJECT}:${VERSION}
	podman push ${REMOTE}/${REPO}/${PROJECT}:${VERSION}

deploy:
	podman-compose -f compose.yaml up --remove-orphans -d

undeploy:
	podman-compose -f compose.yaml down

run:
	go run cmd/webcache/main.go

test:
	go test -v -race ./...