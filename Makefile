VERSION=0.1.1
PROJECT=unsplash-api
REPO=alvidir

PROTOC_GEN_TS_PATH="./node_modules/.bin/protoc-gen-ts"
GRPC_TOOLS_NODE_PROTOC="./node_modules/.bin/grpc_tools_node_protoc"
GRPC_TOOLS_NODE_PROTOC_PLUGIN="./node_modules/.bin/grpc_tools_node_protoc_plugin"

PROTO_DIR="./proto"
OUT_JS_DIR="./src/proto"
OUT_TS_DIR="./src/proto"

build:
	docker build --rm \
	-t ${REPO}/${PROJECT}:${VERSION} \
	-f ./docker/api/dockerfile .
	
	#docker build --rm \
	#-t ${REPO}/${PROJECT}:${VERSION}-envoy \
	#-f ./docker/envoy/dockerfile .

stop:
	docker stop ${PROJECT}
	docker rm ${PROJECT}

	docker stop ${PROJECT}-envoy
	docker rm ${PROJECT}-envoy

logs:
	docker logs -f ${PROJECT}

protoc:
	${GRPC_TOOLS_NODE_PROTOC} \
    	--js_out=import_style=commonjs,binary:${OUT_JS_DIR} \
    	--grpc_out=${OUT_JS_DIR} \
    	--plugin=protoc-gen-grpc=${GRPC_TOOLS_NODE_PROTOC_PLUGIN} \
    	-I ${PROTO_DIR} \
    	${PROTO_DIR}/*.proto
	
	${GRPC_TOOLS_NODE_PROTOC} \
		--ts_out=service=grpc-node:${OUT_TS_DIR} \
    	--plugin=protoc-gen-ts=${PROTOC_GEN_TS_PATH} \
		--plugin="protoc-gen-grpc=${GRPC_TOOLS_NODE_PROTOC_PLUGIN}" \
    	-I ${PROTO_DIR} \
    	${PROTO_DIR}/*.proto

	protoc \
	-I=. ${PROTO_DIR}/*.proto \
	--plugin ./node_modules/.bin/protoc-gen-grpc-web \
	--js_out=import_style=commonjs:./src \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src