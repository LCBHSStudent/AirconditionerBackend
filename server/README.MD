# 分布式温控计费系统
c/s模型，tcp通讯，使用json传输数据，数据库采用mysql

### 启动方式
1. cd server/main
2. go run main.go
3. 或者go build 得到可执行文件之后再运行

## 服务端实现：

### 一，数据结构/数据库表字段设计
```
//1，用户结构
type User struct {
	RoomNum   int    `json:"room_num"`  // 房间号
	Privilege int    `json:"privilege"` // 用户等级，0普通用户，1空调管理员，2酒店经理；不同等级的用户支持对数据库的不同操作
	Password  string `json:"password"`  // 密码
	CheckIn   int64  `json:"check_in"`  //入住时间
	CheckOut  int64  `json:"check_out"` //退房时间
}
```

```
// 2，空调状态结构
type AirConditioner struct {
	Number      int     `json:"number"`      // 空调编号，默认和房间编号相同
	Power       int     `json:"power"`       // 电源开关：0关 1开
	Mode        int     `json:"mode"`        // 模式
	WindLevel   int     `json:"wind_level"`  // 风速
	Temperature float64 `json:"temperature"` // 温度
}
```

```
// 3，房间状态结构体，记录房间号、开始时间、结束时间，以及这段时间内的温度和费用、耗电量
type RoomState struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	StartTime   int64   `json:"start_time"`  // 开始时间
	EndTime     int64   `json:"end_time"`    // 结束时间
	Power       float64 `json:"power"`       // 耗电量
	Cost        float64 `json:"cost"`        // 费用
	Temperature float64 `json:"temperature"` // 温度
}
```

### 二，支持的数据库操作方法

1. 用户方法
```
// 添加用户，需要传入一个user指针，返回数据库实际操作的行数和一个错误信息
func AddUser(user *User) (int64, error) {}

// 通过房间号查询用户，需要传入一个房间号，返回一个user实例
func QueryUserWithRoomNum(roomNum int) User {}

// 更新用户的入住时间，需要传入房间号和用户入住时间
func UpdateUserCheckIn(roomNum int, checkIn int64) (int64, error) {}

// 更新用户的退房时间，需要传入房间号和退房时间
func UpdateUserCheckOut(roomNum int, checkOut int64) (int64, error) {}

// 查询所有用户，该方法返回一个user列表
func QueryAllUsers() (users []User,err error) {}

```
2.空调支持的方法
```
// 往数据库中增加一条空调记录
func AddAirConditioner(air *AirConditioner) (int64, error) {}

// 通过空调编号查询空调状态
func QueryAirWithNumber(number int) (air AirConditioner) {}

// 修改空调状态数据
func UpdateAirConditioner(air *AirConditioner) (int64, error) {}

// 查询所有空调
func QueryAllAirConditioners() (airs []AirConditioner, err error) {}
```
3.房间状态支持的方法
```
// 添加一条房间的温度、耗电量、费用记录
func AddRoomState(roomState *RoomState) (int64, error) {}

// 查询指定时间段内的温度、耗电量、费用记录
func QueryRoomStates(roomNum int, startTime, endTime int64) (roomStates []RoomState, err error) {}

// 删除某个房间指定时间段内的温度、耗电量、费用记录（删除太老的数据）
func DelRoomStatesByRoomNum(roomNum int, startTime, endTime int64) (int64, error) {}
```

### 二，通讯消息结构和消息类型设计

1，先定义外层消息结构，再根据具体消息类型定义详细的消息结构
```
// 定义消息结构体
type Message struct {
	Length int `json:"length"` // 消息长度，用于验证包是否缺失
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 具体的消息数据
}
```

2，具体的消息类型，每一个客户端发送的消息类型都对应一个服务端返回的消息类型
```
const (
    // 用户类型
	TypeUserRegister    = "UserRegister"
	TypeUserRegisterRes = "UserRegisterRes"
	TypeUserLogin       = "UserLogin"
	TypeUserLoginRes    = "UserLoginRes"
	TypeUserModify      = "UserModify"
	TypeUserModifyRes   = "UserModifyRes"
	TypeUserQuery       = "UserQuery"
	TypeUserQueryRes    = "UserQueryRes"
    
    // 空调类型
	TypeAirModify    = "AirModify"
	TypeAirModifyRes = "AirModifyRes"
	TypeAirQuery     = "AirQuery"
	TypeAirQueryRes  = "AirQueryRes"
    
    // 房间状态类型
	TypeRoomStateAdd       = "RoomStateAdd"
	TypeRoomStateAddRes    = "RoomStateAddRes"
	TypeRoomStateQuery     = "RoomStateQuery"
	TypeRoomStateQueryRes  = "RoomStateQueryRes"
	TypeRoomStateDelete    = "RoomStateDelete"
	TypeRoomStateDeleteRes = "RoomStateDeleteRes"
)
```

3，具体的消息结构
```
// 用户注册消息结构体
type UserRegister struct {
	RoomNum   int    `json:"room_num"`
	Privilege int    `json:"privilege"`
	Password  string `json:"password"`
}

// 用户注册返回的消息
type UserRegisterRes struct {
	Code  int    `json:"code"`  // 状态码 400 表示用户已被注册，200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 修改空调状态
type AirModify struct {
	model.AirConditioner `json:"air_conditioner"`
}

// 修改空调状态返回的消息
type AirModifyRes struct {
	Code  int    `json:"code"`  // 状态码
	Error string `json:"error"` // 返回错误信息
}
。。。
以下还有很多，就不一一列出了，具体的消息结构是和消息类型一一对应的
```

### 三，处理服务端和客户端的通信问题

1，主函数，初始化连接池、等待与客户端连接，每当侦听到客户端的连接请求，就启动一个协程，维持和客户端的连接
```
func main() {
    // 初始化数据库
	db.InitMysql(dbName, dsn, maxOpen, maxIdle, maxLifeTime)
    // 在8888端口
	listen,err := net.Listen("tcp", "0.0.0.0:8888")
    // 循环侦听客户端的连接请求，一旦侦听到，就创建一个协程，和客户端保持连接
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
	processor := &Processor.MainProcessor{Conn: conn}
	err := processor.Process()
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
	Conn    net.Conn
	RoomNum int
}
// 用户注册，传入一个msg指针，返回err信息
func (this *UserProcessor) ProcessUserRegister(msg *message.Message) (err error) {}

// 用户登录
func (this *UserProcessor) ProcessUserLogin(msg *message.Message) (err error) {}

// 修改用户的入住时间或者退房时间
func (this *UserProcessor) ProcessUserModify(msg *message.Message) (err error) {}

// 用户查询 如果传入的房间号为0，则查询所有用户信息，否则按房间号查询
func (this *UserProcessor) ProcessUserQuery(msg *message.Message) (err error) {}

// 2.空调处理器和房间状态处理器和用户处理器类似，这里就不列举了
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

2，sql中实现了数据库的初始化工作，以及初始的数据库表的创建工作和操作数据库、查询数据库的方法
```
// 初始化mysql、创建连接池
func InitMysql(name, dsn string, maxOpen, maxIdle int, maxLifeTime time.Duration) {
	db1, err := sql.Open(name, dsn)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		db = db1
		db.SetMaxOpenConns(maxOpen)
		db.SetMaxIdleConns(maxIdle)
		db.SetConnMaxLifetime(maxLifeTime)
		CreateTableWithUser()
		CreateTableWithAirConditioner()
		CreateTableWithRoomStates()
	}
}

// 操作数据库
func ModifyDB(sql string, args ...interface{}) (int64, error) {
	result, err := db.Exec(sql, args...)
	count, err := result.RowsAffected()
	return count, err
}

// 查询单行数据
func QueryRowDB(sql string) *sql.Row {
	return db.QueryRow(sql)
}

// 查询多行数据
func QueryRowsDB(sql string) (*sql.Rows, error) {
	return db.Query(sql)
}

// 创建用户表
func CreateTableWithUser() {}

// 创建空调状态表
func CreateTableWithAirConditioner() {}

// 创建房间状态表，存储温度、耗电量和费用
func CreateTableWithRoomStates() {}
```
