FROM golang:alpine AS build-env
RUN mkdir /go/src/app && apk update && apk add git
ADD main.go /go/src/app/
RUN mkdir /go/.kube
WORKDIR /go/src/app
RUN go mod init
RUN go mod tidy
RUN go get k8s.io/client-go@latest
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

FROM scratch
WORKDIR /app
ADD /var/lib/jenkins/.kube/config /var/lib/jenkins/workspace/.kube/config
COPY --from=build-env /go/src/app/app .
ENTRYPOINT [ "./app" ]
