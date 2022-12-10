FROM golang:1.19.4-alpine3.17 AS src
RUN apk update && apk add --no-cache git
ENV CGO_ENABLED=0
WORKDIR /workspace
#COPY go.mod go.sum ./
#RUN go mod download
COPY . .

FROM src AS test
CMD ["go", "test", "-v", "./..."]

FROM src AS sleeper
ENTRYPOINT ["sleeper/sleeper"]
RUN cd sleeper && GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o sleeper
