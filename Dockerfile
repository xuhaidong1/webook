# 基础镜像  docker build -t xuhaidong/webook:v0.0.1 .
FROM ubuntu:latest
#把编译后的文件打包进镜像中去 放在工作目录/app下
COPY webook /app/webook
WORKDIR /app
#执行命令
ENTRYPOINT ["/app/webook"]