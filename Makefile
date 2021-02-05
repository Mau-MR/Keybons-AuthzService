#protobuf
gen:
	protoc --proto_path=proto proto/*.proto -I. --go_out=plugins=grpc:/home/mawi/gocode/src
clean:
	rm pb/*.go
#running code locally
server:
	go run cmd/server/main.go -port 9090
client:
	go run cmd/client/main.go -address 0.0.0.0:8080
test:
	go test -cover -race ./...
########## container related ###############
container:
	docker build ./docker/envoy
	docker run -d -p 7000:1000 ./docker/envoy
envoyDiagnose:
	docker run --rm \
		-v $$(pwd)/docker/envoy/envoy.yaml:/envoy.yaml \
		envoyproxy/envoy:v1.17-latest \
		--mode validate \
		-c envoy.yaml
serverup:
	docker build -t authz -f $$(pwd)/docker/server/Dockerfile . && docker run -d -p 9090:9090 authz

fullServer:
	docker-compose -f $$(pwd)/docker/docker-compose.yml up -d



.PHONY: gen clean server test client
