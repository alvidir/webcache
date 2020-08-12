export VERSION=0.1.0
export PROJECT=unsplash-api
export REPO=alvidir

export PROTOC_GEN_TS_PATH="./node_modules/.bin/protoc-gen-ts"
export PROTO_PATH="./src/proto"
export OUT_DIR="."

build:
	docker build --rm \
	-t ${REPO}/${PROJECT}:${VERSION} -f ./dockerfile .

run:
	docker run -p 3001:3001 \
	--name=${PROJECT} \
	--restart always \
	--detach ${REPO}/${PROJECT}:${VERSION}

stop:
	docker stop ${PROJECT}
	docker rm ${PROJECT}

logs:
	docker logs -f ${PROJECT}

protoc:
	protoc \
		--plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
		--js_out="import_style=commonjs,binary:${OUT_DIR}" \
		--ts_out="service=grpc-web:${OUT_DIR}" \
		${PROTO_PATH}/*.proto

clear-proto:
	find ./src/proto -name \*.js -type f -delete
	find ./src/proto -name \*.ts -type f -delete