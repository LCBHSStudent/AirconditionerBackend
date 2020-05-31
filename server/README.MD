# 分布式温控计费系统
c/s模型，tcp通讯，使用json传输数据，数据库采用mysql

### 架构设计
* 通讯方式采取tcp通讯
* 数据传输采取json传输
* 使用gorm操作数据库
* 使用viper进行配置管理

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

## 服务端实现：

### 一，数据结构/数据库表字段设计
```
// 用户结构
type User struct {
	Id       int    `json:"id"`
	RoomNum  int    `json:"room_num"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	CheckIn  int64  `json:"check_in"`  //入住时间
	CheckOut int64  `json:"check_out"` //退房时间
}
```

```
// 空调数据结构
type AirConditioner struct {
	Id          int     `json:"id"`          // 空调编号
	RoomNum     int     `json:"room_num"`    // 空调所在房间号，一个房间可能有多个空调
	Power       int     `json:"power"`       // 电源开关：0关 1开
	Mode        int     `json:"mode"`        // 模式
	WindLevel   int     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 温度
}
```

```
// 房间状态结构体，记录房间号、开始时间、结束时间，以及这段时间内的温度和耗电量
type RoomState struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	StartTime   int64   `json:"start_time"`  // 开始时间
	EndTime     int64   `json:"end_time"`    // 结束时间
	Power       float64 `json:"power"`       // 耗电量
	Temperature float64 `json:"temperature"` // 温度
}
```

```
// 费用结构体，记录房间号、开始时间、结束时间，以及这段时间内的费用
type Fee struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	StartTime   int64   `json:"start_time"`  // 开始时间
	EndTime     int64   `json:"end_time"`    // 结束时间
	Cost        float64 `json:"cost"`        // 费用
}
```

### 二，支持的数据库操作方法

1. 用户方法
```
// 用户数据库操作相关接口
type UserRepository interface {
	FindById(int) (*model.User, error)                        // 根据Id查找用户
	Create(*model.User) error                                 // 创建用户
	Update(*model.User) (*model.User, error)                  // 更新用户
	FindByField(string, string, string) ([]model.User, error) // 条件查询
	FindAllUsers() (users []model.User, err error)            // 查询所有用户
}
```
2.空调支持的方法
```
// 空调数据库操作相关接口
type AirConditionerRepository interface {
	FindByID(int) (*model.AirConditioner, error)                       // 根据id查找空调
	Create(*model.AirConditioner) error                                // 创建空调状态
	Update(*model.AirConditioner) error                                // 更新空调状态
	FindByField(string, string, string) (*model.AirConditioner, error) // 条件查询
	FindAll() ([]model.AirConditioner, error)                          // 查询所有空调
	FindByRoom(int) ([]model.AirConditioner, error)                    // 通过房间号查询
}
```
3.房间状态支持的方法
```
// 房间状态数据库操作相关接口
type RoomStateRepository interface {
	AddRoomState(roomState *model.RoomState) error                                    // 添加一条房间状态信息
	QueryRoomStates(roomNum int, startTime, endTime int64) ([]model.RoomState, error) // 查询指定时间段内的房间状态，返回房间状态数组
	DelRoomStates(roomNum int, startTime, endTime int64) (int64, error)               // 删除某个房间指定时间段内的状态记录（删除太老的数据），返回删除条数
}
```
4.账单
```
// 账单数据库操作相关接口
type FeeRepository interface {
	AddFee(roomState *model.Fee) error                                    // 添加一条账单信息
	QueryFees(roomNum int, startTime, endTime int64) ([]model.Fee, error) // 查询指定时间段内的账单，返回账单数组
	DelFees(roomNum int, startTime, endTime int64) (int64, error)         // 删除某个房间指定时间段内的账单记录（删除太老的数据），返回删除条数
}
```

### 二，通讯消息结构和消息类型设计

1，先定义外层消息结构，再根据具体消息类型定义详细的消息结构
```
// 定义消息结构体
type Message struct {
	Length int    `json:"length"` // 消息长度，用于验证包是否缺失
	Type   string `json:"type"`   // 消息类型
	Data   string `json:"data"`   // 消息
}
```

2，具体的消息类型，每一个客户端发送的消息类型都对应一个服务端返回的消息类型
```
// 消息类型
const (
	TypeUserRegister    = "UserRegister"
	TypeUserRegisterRes = "UserRegisterRes"
	TypeUserLogin       = "UserLogin"
	TypeUserLoginRes    = "UserLoginRes"
	TypeUserFindById    = "UserFindById"
	TypeUserFindByIdRes = "UserFindByIdRes"
	TypeUserFindAll     = "UserFindAll"
	TypeUserFindAllRes  = "UserFindAllRes"
	TypeUserUpdate      = "UserUpdate"
	TypeUserUpdateRes   = "UserUpdate"

	TypeAirConditionerFindById      = "AirConditionerFindById"
	TypeAirConditionerFindByIdRes   = "AirConditionerFindByIdRes"
	TypeAirConditionerFindByRoom    = "AirConditionerFindByRoom"
	TypeAirConditionerFindByRoomRes = "AirConditionerFindByRoomRes"
	TypeAirConditionerFindAll       = "AirConditionerFindAll"
	TypeAirConditionerFindAllRes    = "AirConditionerFindAllRes"
	TypeAirConditionerCreate        = "AirConditionerCreate"
	TypeAirConditionerCreateRes     = "AirConditionerCreateRes"
	TypeAirConditionerUpdate        = "AirConditionerUpdate"
	TypeAirConditionerUpdateRes     = "AirConditionerUpdateRes"

	TypeRoomStateAdd       = "RoomStateAdd"
	TypeRoomStateAddRes    = "RoomStateAddRes"
	TypeRoomStateQuery     = "RoomStateQuery"
	TypeRoomStateQueryRes  = "RoomStateQueryRes"
	TypeRoomStateDelete    = "RoomStateDelete"
	TypeRoomStateDeleteRes = "RoomStateDeleteRes"

	TypeFeeAdd       = "FeeAdd"
	TypeFeeAddRes    = "FeeAddRes"
	TypeFeeQuery     = "FeeQuery"
	TypeFeeQueryRes  = "FeeQueryRes"
	TypeFeeDelete    = "FeeDelete"
	TypeFeeDeleteRes = "FeeDeleteRes"
)
```

3，具体的消息结构
```
// 用户注册消息结构体
type UserRegister struct {
	RoomNum  int    `json:"room_num"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户注册返回的消息
type UserRegisterRes struct {
	Code int    `json:"code"`  // 状态码 400 表示用户已被注册，200表示注册成功
	Msg  string `json:"error"` // 返回错误信息
}

// 用户登录消息结构体
type UserLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户登录返回的消息
type UserLoginRes struct {
	Code int    `json:"code"`  // 状态码
	Msg  string `json:"error"` // 返回错误信息
}
。。。
以下还有很多，就不一一列出了，具体的消息结构是和消息类型一一对应的
```

### 三，处理服务端和客户端的通信问题

1，主函数，初始化连接池、等待与客户端连接，每当侦听到客户端的连接请求，就启动一个协程，维持和客户端的连接
```
// 主函数，初始化连接池、与客户端连接、启动协程
func main() {

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("connection error : %v \n", err)
	}
	defer db.Close()

	fmt.Println("服务器在8888端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("net.listen err=", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		fmt.Println("和客户端连接成功...")
		go mainProcess(conn)
	}
}

// 主控程序，处理和客户端的通讯
func mainProcess(conn net.Conn) {
	defer conn.Close()
	mainProcessor := &processor.MainProcessor{Conn: conn}
	err := mainProcessor.Process()
	if err != nil {
		fmt.Println("通讯协程错误，err=", err)
		return
	}
}
```

2,主调度器，用于根据客户端发送的消息类型来分配给具体的处理器进行处理
```
// 主调度器结构，维持一个从主函数中获取的连接
type MainProcessor struct {
	Conn net.Conn
}

// Process用于监听并读取客户端发来的消息
func (this *MainProcessor) Process() (err error) {
	conn := this.Conn
	for {
		tf := &utils.Transfer{Conn: conn,}
		msg, err := tf.ReadPkg()
		err = this.serverProcessMsg(&msg)
    }
}

// 功能：根据客户端发送的消息种类不同，决定调用哪个函数来处理
func (this *MainProcessor) serverProcessMsg(msg *message.Message) (err error) {
	conn := this.Conn
	switch msg.Type {...}
}
```
3，各个从处理器，各司其职

```
// 1，用户处理器，处理和用户相关的消息
type UserProcessor struct {
	Conn net.Conn
	Orm  *repository.UserOrm
}

// 用户注册
func (up *UserProcessor) Register(msg *message.Message) (err error) {}

// 用户登录
func (up *UserProcessor) Login(msg *message.Message) (err error) {}

// 用户查询
func (up *UserProcessor) FindById(msg *message.Message) (err error) 

// 查询所有用户
func (up *UserProcessor) FindAll(msg *message.Message) (err error)

// 2.其余处理器类似，这里就不列举了
```

### 小工具

1,utils中实现了具体的读写方法，其中读取消息的时候会先判断消息长度是否和结构体内的消息长度一致，如果不一致，认为消息损坏，直接丢弃
```
// 定义一个传输结构体，以便绑定读写方法到一个连接上
type Transfer struct {
	Conn   net.Conn
	Buffer [4096]byte
}
// 从连接中读取信息，读取成功则返回一个msg
func (this *Transfer) ReadPkg() (msg message.Message, err error) {}

// 向一个连接中写入消息，需要传入一个字节切片
func (this *Transfer) WritePkg(data []byte) (err error) {}
```

