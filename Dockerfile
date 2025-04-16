
FROM golang:1.23-alpine AS build

WORKDIR /go/src

COPY go.mod ./
RUN go mod download
COPY . .
RUN go install -v ./...

FROM golang:1.23-alpine
COPY --from=build /go/bin/crowemiwebhooks /go/crowemiwebhooks

EXPOSE 8003

CMD ["/go/crowemiwebhooks"]