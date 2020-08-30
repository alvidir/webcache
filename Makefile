VERSION=0.1.1
PROJECT=unsplash-api
REPO=alvidir

build:
	docker build --rm \
	-t ${REPO}/${PROJECT}:${VERSION} \
	-f ./docker/api/dockerfile .

run:
	docker run -it \
	--publish 3030:3030 \
	--name ${PROJECT} \
	${REPO}/${PROJECT}:${VERSION}

stop:
	docker stop ${PROJECT}
	docker rm ${PROJECT}