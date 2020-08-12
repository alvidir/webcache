VERSION=0.1.0
PROJECT=unsplash-api
REPO=alvidir

PROTOC_GEN_TS_PATH="./node_modules/.bin/protoc-gen-ts"
GRPC_TOOLS_NODE_PROTOC_PLUGIN="./node_modules/.bin/grpc_tools_node_protoc_plugin"
GRPC_TOOLS_NODE_PROTOC="./node_modules/.bin/grpc_tools_node_protoc"

PROTO_DIR="./proto"
OUT_DIR="./src/model"

build:
	docker build --rm \
	-t ${REPO}/${PROJECT}:${VERSION} -f ./dockerfile .

run:
	docker run -p 3001:3001 \
	--name=${PROJECT} \
	--restart always \
	-it ${REPO}/${PROJECT}:${VERSION}

stop:
	docker stop ${PROJECT}
	docker rm ${PROJECT}

logs:
	docker logs -f ${PROJECT}

protoc:
	protoc \
		--plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
		--plugin=protoc-gen-grpc=${GRPC_TOOLS_NODE_PROTOC_PLUGIN} \
		--js_out="import_style=commonjs,binary:${OUT_DIR}" \
		--ts_out="service=grpc-web:${OUT_DIR}" \
		--grpc_out=${OUT_DIR} \
		--proto_path ${PROTO_DIR} ${PROTO_DIR}/*.proto

	#${GRPC_TOOLS_NODE_PROTOC} \
  	#	--plugin="protoc-gen-ts=${PROTOC_GEN_TS_PATH}" \
  	#	--grpc_out="${OUT_DIR}" \
  	#	--js_out="import_style=commonjs,binary:${OUT_DIR}" \
  	#	--ts_out="${OUT_DIR}" \
  	#	--proto_path ${PROTO_DIR} ${PROTO_DIR}/*.proto

clear-proto:
	find ./src/model -name \*pb.js -type f -delete
	find ./src/model -name \*pb.d.ts -type f -delete