NAME:=whdbg
REGISTRY:=zushco
WEBSOCKET:="wss://whdbg.dev/"
all:
	go build -o bin/$(NAME) ./...

linux:
		GOOS=linux GOARCH=amd64 go build -o bin/$(NAME) ./...

.PHONY: web
web:
	(cd web && yarn build)

.PHONY: start
start:
	(cd web && yarn start)

image:
	docker buildx build --platform linux/amd64 --build-arg prod_websocket=$(WEBSOCKET) -t $(NAME) .

push:
	docker tag $(NAME) registry.digitalocean.com/$(REGISTRY)/$(NAME)

deploy:
	docker buildx build --push \
		--platform linux/amd64 \
		--tag registry.digitalocean.com/zushco/$(NAME):latest  .

run: web
	go run *.go
