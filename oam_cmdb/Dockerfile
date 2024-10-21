# 生成镜像前先修改好配置文件
FROM library/golang:1.18

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR $GOPATH/src/OAM
COPY . $GOPATH/src/OAM
RUN --mount=type=cache,mode=0777,id=go-mod,target=/go/pkg/mod \
	go build -ldflags '-d -w -s'

EXPOSE 8848
ENTRYPOINT ["./OAM"]