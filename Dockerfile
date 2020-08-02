# 基础镜像
FROM golang:latest as build
ENV GOPROXY=https://goproxy.cn,direct
ENV GOPRIVATE=github.com/zengyu2020
WORKDIR /go/src/app
COPY . .
RUN  git config --global url."http://github.com/zengyu2020".insteadOf "git@gitlab.com:zengyu2020" && git config --global http.extraheader "PRIVATE-TOKEN: b4486da5cf22d0b71046d25d194b31be839d830f" && go build -o im-static main.go

FROM alpine:latest
WORKDIR /program
COPY --from=build /go/src/app/im-static .
VOLUME ["/program/upload"]
RUN mkdir /lib64 && \
 ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
 echo "https://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories &&\
 apk add -U tzdata && \
 ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
 echo "Asia/Shanghai" > /etc/timezone && \
 chmod +x im-static
ENTRYPOINT ["./im-static"]
EXPOSE 8091