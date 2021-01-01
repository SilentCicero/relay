FROM golang AS build
RUN echo 'nobody:x:65534:65534:Nobody:/:' > /passwd
WORKDIR $GOPATH/src/gitlab.com/jonas.jasas/httprelay
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X gitlab.com/jonas.jasas/httprelay/pkg/server.Version=`git describe --tag`" -o /httprelay ./cmd/...

FROM scratch
COPY --from=build /passwd /etc/passwd
COPY --from=build /httprelay .
USER nobody
EXPOSE 8080
ENTRYPOINT ["/httprelay"]