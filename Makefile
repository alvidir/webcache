export VERSION=0.1.0
export PROJECT=unsplash-api
export REPO=alvidir

build:
	docker build -t ${REPO}/${PROJECT}:${VERSION} -f ./dockerfile .

run:
	docker run -p 3001:3001 \
	--name=${PROJECT} \
	--restart always \
	--detach ${REPO}/${PROJECT}:${VERSION}

stop:
	docker stop ${PROJECT}
	docker rm ${PROJECT}