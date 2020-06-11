

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
	RoomNum         int     `json:"room_num"`         // 空调所在房间号，默认一个房间一个空调(中央空调形式)
	Power           int     `json:"power"`            // 电源开关：0关 1开
	Mode            int     `json:"mode"`             // 模式
	WindLevel       int     `json:"wind_level"`       // 风速
	Temperature     float64 `json:"temperature"`      // 温度
	RoomTemperature float64 `json:"room_temperature"` // 室温
	TotalPower      float64 `json:"total_power"`      // 该次入住的总耗电量
	StartWind       []int64 `json:"start_wind"`       // 开始送风时间，时间戳格式
	StopWind        []int64 `json:"stop_wind"`        // 停止送风时间
	OpenTime        []int64 `json:"open_time"`        // 开机时间，数组
	CloseTime       []int64 `json:"close_time"`       // 关机时间，数组
	SetParamNum     int     `json:"set_param_num"`    // 调整次数，用于报表展示
}
```


```
// 费用结构体，记录房间号、费用
type Fee struct {
	RoomNum     int     `json:"room_num"`    // 房间号
	Cost        float64 `json:"cost"`        // 费用
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
	TypeNormalRes = "NormalRes"

	TypeUserRegister    = "UserRegister"
	TypeUserLogin       = "UserLogin"
	TypeUserFindById    = "UserFindById"
	TypeUserFindByIdRes = "UserFindByIdRes"
	TypeUserFindAll     = "UserFindAll"
	TypeUserFindAllRes  = "UserFindAllRes"
	TypeUserUpdate      = "UserUpdate"

	TypeAirConditionerFindByRoom    = "AirConditionerFindByRoom"
	TypeAirConditionerFindByRoomRes = "AirConditionerFindByRoomRes"
	TypeAirConditionerFindAll       = "AirConditionerFindAll"
	TypeAirConditionerFindAllRes    = "AirConditionerFindAllRes"
	TypeAirConditionerCreate        = "AirConditionerCreate"
	TypeAirConditionerUpdate        = "AirConditionerUpdate"
	TypeAirConditionerOn            = "AirConditionerOn"
	TypeAirConditionerOff           = "AirConditionerOff"
	TypeAirConditionerSetParam      = "AirConditionerSetParam"
	TypeAirConditionerStopWind      = "AirConditionerStopWind"
	TypeGetReport                   = "GetReport"
	TypeGetReportRes                = "GetReportRes"

	TypeFeeAdd       = "FeeAdd"
	TypeFeeQuery     = "FeeQuery"
	TypeFeeQueryRes  = "FeeQueryRes"
	TypeFeeDelete    = "FeeDelete"
)
```

3，具体的消息结构
```
// 普通的消息返回格式，只包含状态码和msg，如需包含额外数据则另外定义返回格式
type NormalRes struct {
	Code int    `json:"code"` // 状态码
	Msg  string `json:"msg"`  // 返回信息
}

// 用户注册消息结构体
type UserRegister struct {
	RoomNum  int    `json:"room_num"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// 用户登录消息结构体
type UserLogin struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserFindById struct {
	Id int `json:"id"`
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

