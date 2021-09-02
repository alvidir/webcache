VERSION=0.1.0
PROJECT=webcache
REPO=alvidir

build:
	podman build -t ${REPO}/${PROJECT}:${VERSION}-nginx -f ./docker/nginx/dockerfile .
	podman build -t ${REPO}/${PROJECT}:${VERSION}-server -f ./docker/webcache/dockerfile .

deploy:
	podman-compose -f docker-compose.yaml up --remove-orphans
	# delete -d in order to see output logs

undeploy:
	podman-compose -f docker-compose.yaml down