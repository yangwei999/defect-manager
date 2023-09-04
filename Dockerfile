FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/defect-manager
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -a -o defect-manager .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 defect && \
    useradd -u 1000 -g defect -s /sbin/nologin -m defect

USER defect

COPY  --chown=defect --from=BUILDER /go/src/github.com/opensourceways/defect-manager/defect-manager /opt/app/defect-manager

ENTRYPOINT ["/opt/app/defect-manager"]