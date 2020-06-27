## 《分布式温控系统》静态结构设计

​		本文档是基于动态结构设计描述的系统静态结构，包括了客户端与服务端的类与方法，消息交互格式以及系统级静态结构。

#### 一、客户端

1. 顾客

   1.1方法

   1.1.1  void AirConditionOn()

​	方法描述：空调开机

| 参数       | 类型  | 说明               |
| ---------- | ----- | ------------------ |
| roomNum    | int   | 房间号             |
| targetTemp | float | 目标温度           |
| pattern    | int   | 模式（制热/制冷）  |
| windspeed  | int   | 风速（微/中/大风） |

​	1.1.2 void AirConditionOff()

​	方法描述：空调关机

| 参数    | 类型 | 说明   |
| ------- | ---- | ------ |
| roomNum | int  | 房间号 |

​	1.1.3 void RequestService()

​	方法描述：向服务器请求（设置模式/目标温度/风速）

| 参数       | 类型  | 说明               |
| ---------- | ----- | ------------------ |
| roomNum    | int   | 房间号             |
| targetTemp | float | 目标温度           |
| pattern    | int   | 模式（制热/制冷）  |
| windspeed  | int   | 风速（微/中/大风） |

​	1.1.4  void stopService()

​	方法描述：当达到目标温度，自动发送停止送风请求。

| 参数    | 类型 | 说明   |
| ------- | ---- | ------ |
| roomNum | int  | 房间号 |

​	1.1.5 void sendTotalPower()

​	方法描述：客户端定时发送耗电量

| 参数       | 类型  | 说明   |
| ---------- | ----- | ------ |
| roomNum    | int   | 房间号 |
| totalPower | float | 耗电量 |



2. 管理员

   2.1 方法

   void getMechineState() 

   方法描述：监控房间空调，定时向服务器请求获取各个房间的空调使用情况，接收到服务器的回应并展示收到的数据。



3. 经理

   3.1 方法

   ​	3.1.1 void getStatement()

   方法描述：向服务器发送查看酒店的格式化统计报表（日报表）的请求，接受来自服务器的格式化统计报表并输出到显示界面。

   ​	3.1.2 void printStatement()

   方法描述：打印格式化统计报表到文件中。

   

4. 前台

   4.1 方法

   ​	4.1.1  void getAccount()

   方法描述：向服务器请求获取房间号为roomNum的空调使用账单。接收服务器的回送并展示到页面。

   | 参数    | 类型 | 说明   |
   | ------- | ---- | ------ |
   | roomNum | int  | 房间号 |

   ​	4.1.2 void getDetailAccount()

   方法描述：向服务器请求获取房间号为roomID的空调详细账单。接收服务器的回送并展示到页面。

   | 参数    | 类型 | 说明   |
   | ------- | ---- | ------ |
   | roomNum | int  | 房间号 |

   ​	4.1.3 void printAccount()

   方法描述：打印房间号为roomNum的空调的账单到文件。

   | 参数    | 类型 | 说明   |
   | ------- | ---- | ------ |
   | roomNum | int  | 房间号 |

   ​	4.1.4 void printDetailAccount()

   方法描述：打印房间号为roomNum的空调的详细账单到文件。

   | 参数    | 类型 | 说明   |
   | ------- | ---- | ------ |
   | roomNum | int  | 房间号 |

#### 二、服务端

1. 数据结构/数据库表字段设计

​    1.1 AirConditioner——空调结构体

​    1.1.1 属性

| 参数        | 类型    | 说明                         |
| ----------- | ------- | ---------------------------- |
| Id          | int     | 空调编号                     |
| RoomNum     | int     | 空调房间号                   |
| Power       | int     | 1开机/0关机                  |
| Mode        | int     | 模式                         |
| WindLevel   | int     | 风速                         |
| Temperature | float64 | 温度                         |
| totalPower  | float64 | 总耗电量                     |
| totalCost   | float64 | 总花费                       |
| startWind   | []int64 | 开始送风时间                 |
| endWind     | []int64 | 停止送风时间                 |
| openTime    | []int64 | 开机时间                     |
| closeTime   | []int64 | 关机时间                     |
| setNum      | int     | 调整空调参数的次数（报表用） |



   1.1.2 方法

| 方法                                      | 说明           |
| ----------------------------------------- | -------------- |
| Update()                                  | 更新空调状态   |
| FindAll([]AirConditioner,error)           | 查询所有空调   |
| FindByRoomId(int)([]AirConditioner,error) | 通过房间号查询 |
| airconditionSchedule()                    | 空调调度程序   |

 
1.2 Fee——账单结构体

​	1.2.1 属性

| 参数      | 类型    | 说明     |
| --------- | ------- | -------- |
| RoomNum   | int     | 房间号   |
| StartTime | int64   | 开始时间 |
| EndTime   | int64   | 结束时间 |
| Cost      | float64 | 费用     |

​    1.2.2 方法

| 方法              | 说明                                         |
| ----------------- | -------------------------------------------- |
| AddFees()         | 添加一条账单信息                             |
| QueryFees()       | 查询某个房间空调的账单，返回账单数组         |
| QueryDetailFees() | 查询某个房间空调的详细账单，返回详细账单数组 |



#### 三、通信消息结构
注：json中的字段统一采取 下划线分割命名 形式，如 "room_num"、"wind_level" 等

1、顾客

  1.1 开机

客户端发送：

```json
//客户端请求
{
    "type":"AirConditionerOn",
    "data":{
        "room_num": 101,          //房间号
    	"mode": "cold",                //模式（0制冷/1制热）
   	 	"wind_level": "low",          //风速	
    	"temperature": 25.5,      //目标温度
    	"open_time": 1591873406   //开机时间，时间戳格式
    }
}

//服务端返回
{
    "code":200,     //200 ，500 
    "messege":"开机成功！"    //ok ， wait 
}
```

​	1.2 设置参数：模式、风速、温度

```json
//客户端请求
{
    "type":"AirConditionerSetParam",
    "data":{
        "room_num":101, //房间号
    	"mode": "cold",                //模式（0制冷/1制热）
   	 	"wind_level": "high",          //风速
    	"temperature":25.5
    }
}
//服务端返回
{
    "code":200,   //200 / 500 / 404
    "messege":"调整空调参数成功！"  //ok / wait /
}
```

​	1.3 关机

```json
//客户端发送
{
    "type":"AirConditionerOff",
    "data":{
        "room_num":101, //房间号
    	"close_time":1591873406
    }
}
//服务端返回
{
    "code":200,   
    "messege":"关机成功！"  
}
```

​	1.4  定时发送房间（总）耗电量和房间温度

```json
//客户端发送
{
	"type":"SetRoomData",
    "data":{
        "room_num":101, //房间号
        "total_power":12.3, // 总耗电量
        "room_temperature":28 // 房间温度
    }
}
//服务端返回
{
    "code":"200",   
    "messege":"设置空调耗电量和房间温度成功！"  
}
```

​	1.5 停止送风请求

```json
//客户端发送
{
    "type":"AirConditionerStopWind",
    "data":{
      "room_num":101, //房间号
      "end_wind":1591873406,   //停止送风时间戳
    }
}
//服务端返回
{
    "code":200,  
    "messege":"空调停止送风成功！"  
}
```



2、管理员

​    2.1 监控房间空调状态

```json
//客户端请求
{
	"type":"AirConditionerFindAll"
}
//服务端返回
{
	"code":200,   
    "messege":"查询所有空调状态成功！" , 
    "data":[
    	{
    		"room_num":101,
    		"power":"",    //正在送风/等待/关机
            "mode":"",     
            "wind_level":"",
            "temperature":"",
            "total_power":"",  //耗电量
            "fee":"",
            ...,
    	},
    	{
    		"room_num":102,
    		"power":"",    //正在送风/等待/关机
            "mode":"",     
            "wind_level":"",
            "temperature":"",
            "total_power":"",  //耗电量
            "fee":""
    	}
    	...
    ]
}
```

3、经理

​    3.1 查看报表（日/月/年）

```json
//客户端请求
{
 	"type":"GetReport"   
}
//服务端返回
{
    "code":200,  
    "messege":"获取所有空调报表成功！",
    "data":[
      {
        "room_num":"",      //房间号
        "fee":"",           //总费用
        "total_power":"",   //总耗电量
        "close_num": "",    //空调的开关次数
        "set_param_num":"", //空调调整次数
        "used_time":""      //使用空调的时长
      },
      {
          "room_num":"",      //房间号
          "fee":"",           //总费用
          "total_power":"",   //总耗电量
          "close_num": "",    //空调的开关次数
          "set_param_num":"", //空调调整次数
          "used_time":""      //使用空调的时长
       },
        ...
     ]
}
```

4、前台

​    4.1 查看账单

```json
//客户端发送
{
	"type":"FeeQuery",
    "data":
    {
          "room_num":101,         //房间号
    }
}
//服务端返回
{
    "code":200,    
    "messege":"查询账单成功！",  
    "data":
    {
          "room_num":101,         //房间号
          "cost":23.4,            //总费用  
    }     //fee 费用
}
```

​    4.2 查看详单 **待完善

```json
//客户端发送
{
	"type":"GetDetailList",
    "data":
    {
        "room_num":101,         //房间号
    }
}


//服务端返回
{
    "code":200,  
    "messege":"获取详单成功！",  
    "data":{
        "room_num":101,
        "start_wind_list":"",
        "stop_wind_list":"",
        "total_wind_time":"",
        "wind_level":"",
        "total_power":"",
        "fee_rate":"",
        "total_fee":"",
    }
}
```



#### 四、系统级静态结构

![image-20200601214651089](C:\Users\28503\AppData\Roaming\Typora\typora-user-images\image-20200601214651089.png)