FROM alpine:latest
WORKDIR /program
COPY im-static .
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



docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e DRONE_RPC_PROTO=http -e DRONE_RPC_HOST=216.250.106.214:8090 -e DRONE_RPC_SECRET=drone_rpc_secret -e DRONE_RUNNER_CAPACITY=2 -e DRONE_RUNNER_NAME=MyCloudServer -e DRONE_LOGS_TRACE=true -p 8080:3000 --restart=always --name drone_runner drone/drone-runner-docker