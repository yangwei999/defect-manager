FROM golang:1.18.8 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/defect-manager
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o defect-manager .

# copy binary config and utils
FROM alpine:3.14
COPY  --from=BUILDER /go/src/github.com/opensourceways/defect-manager/defect-manager /opt/app/defect-manager

ENTRYPOINT ["/opt/app/defect-manager"]