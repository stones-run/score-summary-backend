

# 编译

## 初始化，生成go.mod
1、go mod init score-summary-backend
## 添加依赖到go.mod
2、go mod tidy
## 无法下载,更换代理
3、go env -w GOPROXY=https://goproxy.cn