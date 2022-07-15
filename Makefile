run_server:
	go run ./server/cmd/tcp/tcp.go
.PHONY: run_server

run_client:
	go run ./client/cmd/tcp/tcp.go
.PNONY: run_client

build_server_image:
	docker build -t f_server -f ./server/Dockerfile .
.PHONY: build_server_image

run_server_docker:
	docker run -p 6677:6677 f_server
.PHONY: run_server_docker

build_client_image:
	docker build -t f_client -f ./client/Dockerfile .
.PHONY: build_client_image

run_client_docker:
	docker run -p 6677:6677 f_client
.PHONY: run_client_docker

