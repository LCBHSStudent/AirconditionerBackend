# AirConditioner

实现一个分布式温控收费系统

server端代码，详情查看server/README.MD

### 启动方式
1. 配置go环境
    * set GO111MODULE=on -- 设置go mod管理依赖
    * set GOPROXY="https://goproxy.cn" -- 设置GOPROXY

2. 启动服务端
    * cd server
    * go run main.go
    * 或者go build 得到可执行文件之后再运行
    
3. 运行客户端进行测试
    * cd client
    * go run main.go