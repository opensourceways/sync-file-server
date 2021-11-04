FROM golang:latest as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
COPY . /go/src/github.com/opensourceways/repo-file-server
RUN cd /go/src/github.com/opensourceways/repo-file-server && CGO_ENABLED=1 go build -v -o ./repo-file-server ./server

# copy binary config and utils
FROM golang:latest
RUN  mkdir -p /opt/app/
# overwrite config yaml
COPY  --from=BUILDER /go/src/github.com/opensourceways/repo-file-server/repo-file-server /opt/app

WORKDIR /opt/app/
ENTRYPOINT ["/opt/app/repo-file-server"]
