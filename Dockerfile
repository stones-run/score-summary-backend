# 使用官方的Go镜像作为基础镜像
FROM golang:1.17

WORKDIR /home/admin/score-summary-backend

# 设置 GOPROXY 环境变量
ENV GOPROXY=https://goproxy.cn

# 将当前目录下的所有文件复制到容器中
COPY . .


# 指定容器运行时要执行的命令
CMD ["sh", "-c", "go build -o score-summary && ./score-summary"]
