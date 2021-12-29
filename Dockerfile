FROM golang:alpine AS build-env
RUN mkdir /go/src/app && apk update && apk add git
ADD main.go /go/src/app/
RUN mkdir /kube
COPY /kube/.kube /kube/.kube
WORKDIR /go/src/app
RUN go mod init
RUN go mod tidy
RUN go get k8s.io/client-go@latest
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

FROM scratch
WORKDIR /app
COPY --from=build-env /go/src/app/app .
ENTRYPOINT [ "./app" ]
