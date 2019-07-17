# build stage
FROM golang:alpine AS build-env
ENV GOPATH=/go
WORKDIR /go
ADD src /go/src
RUN apk add git
RUN go get github.com/shirou/gopsutil/process
RUN CGO_ENABLED=0 go install server

# final stage
FROM alpine
WORKDIR /bin
COPY --from=build-env /go/bin/server /bin/server