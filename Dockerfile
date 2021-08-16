#
# Build Node
FROM docker.io/node:14-alpine3.14 as build-node
RUN apk --no-cache --virtual build-dependencies add \
        python3 \
        make \
        g++

ARG prod_websocket
ENV WDS_SOCKET_HOST=$prod_websocket
WORKDIR /workdir
COPY web/ .
RUN yarn install
RUN yarn build


#
# Build Go
FROM docker.io/golang:1.16-alpine3.14 as build-go
ENV GOPATH ""
RUN go env -w GOPROXY=direct
RUN apk add git

ADD go.mod go.sum ./
RUN go mod download
ADD . .
COPY --from=build-node /workdir/build ./web/build
RUN go build -o /main

FROM docker.io/alpine:3.13 as release
COPY --from=build-go /main /main
ENTRYPOINT [ "/main" ]
