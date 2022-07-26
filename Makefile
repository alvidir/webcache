build:
	podman build -t webcache:latest -f ./container/webcache/containerfile .

deploy:
	podman-compose -f compose.yaml up --remove-orphans

undeploy:
	podman-compose -f compose.yaml down

run:
	go run cmd/webcache/main.go

test:
	go test -v -race ./...