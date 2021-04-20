FROM xiaowang1234567/golang:1.0.0 AS builder
RUN yum -y update && \
    yum install -y libssh2-devel && \
    yum clean all && mkdir -p /libgit2 && \
    git config --global http.postBuffer 2000000000 && \
    git config --global user.name wang1137095129 && \
    git config --global user.email wang1137095129@foxmail.com && \
    git clone https://github.com/libgit2/libgit2 /libgit2 && cd /libgit2 &&\
    mkdir -p build && cd build && cmake .. && \
    cmake --build . --target install && \
    mkdir -p "$GOPATH/github.com/wang1137095129/go-git"
ADD . "$GOPATH/github.com/wang1137095129/go-git"
RUN cd "$GOPATH/github.com/wang1137095129/go-git" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -0 /go_git

FROM bitnami/minideb:stretch
RUN install_packages ca-certificates
COPY --from=builder /go_git /go_git
ENTRYPOINT ["/go_git"]