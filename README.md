# 分布式温控计费系统
c/s模型，tcp通讯，使用json传输数据，数据库采用mysql

### 架构设计
* 通讯方式采取tcp通讯
* 使用json进行数据传输
* 使用gorm操作数据库，包括数据表的创建、增删改查等
* 使用go-viper进行配置管理

#### server端代码，详情查看server/README.MD
#### 接口文档，详情查看InterfaceDoc.MD

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