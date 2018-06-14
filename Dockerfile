FROM golang:1.10.2

RUN mkdir -p /go/src/github.com/jaimemartinez88/api.shopping
WORKDIR /go/src/github.com/jaimemartinez88/api.shopping
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a ./cmd/api.shopping/

FROM scratch
COPY --from=0 /go/src/github.com/jaimemartinez88/api.shopping/api.shopping /api.shopping
COPY --from=0 /go/src/github.com/jaimemartinez88/api.shopping/config.json /config.json
CMD ["/api.shopping","-config","config.json"]