package scheduler

import (
	"fmt"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"gorm.io/gorm"
	"log"
	"math"
	"time"
)

// 空调调度请求队列
var RequestQueue []ScheduleReq

// 当前正在服务的空调队列
var ServingQueue []ScheduleReq

// 当前正在等待服务的空调队列
var WatingQueue	 []ScheduleReq

const (
	// 服务队列最大数量
	MaxServerNum = 3
	// 时间片
	TimeSlice = 120
	// 内层时间片
	TimeGap = 5
)

type WindTime struct {
	RoomNum  	int
	StartWindTime 	[]int64
	StopWindTime 	[]int64
}

var SchedulerDB *gorm.DB

var windTimeMap = make(map[int]WindTime)

func Schedule(){

	for {
		fmt.Println("RequestQueue:", RequestQueue)
		fmt.Println("ServingQueue:", ServingQueue)
		fmt.Println("WatingQueue:", WatingQueue)

		// 处理等待队列和服务队列
		for len(ServingQueue) < MaxServerNum {
			// 如果等待队列为空，就跳出循环了
			if len(RequestQueue) == 0 {
				break
			} else { // 否则就直接往服务队列里调度
				nowReq := RequestQueue[0]
				if flag := checkWindStop(nowReq); flag{
					RequestQueue = RequestQueue[1:]
					continue
				}
				if flag := checkPowerOff(nowReq); flag{
					RequestQueue = RequestQueue[1:]
					continue
				}
				if flag := checkWindChange(nowReq); flag{
					RequestQueue = RequestQueue[1:]
					continue
				}
				ServingQueue = append(ServingQueue, nowReq)
				RequestQueue = RequestQueue[1:]
			}
		}

		// 如果调度请求队列不为空，则逐个处理请求
		if len(RequestQueue) != 0 {
			for i:=0; i < len(RequestQueue); i++ {
				nowReq := RequestQueue[i]

				if flag := checkWindStop(nowReq); flag{
					continue
				}

				if flag := checkPowerOff(nowReq); flag{
					continue
				}

				if flag := checkWindChange(nowReq); flag{
					continue
				}

				// 去服务队列中查找有无比当前请求队列优先级低的，如果有就替换
				for j := 0; j < len(ServingQueue); j++ {
					// 如果新送风请求的风速若高于（高风>中风>低风）正在接受服务的某个送风请求，则将立即服务高风速请求；
					if nowReq.WindLevel > ServingQueue[j].WindLevel{
						// 将原请求放到等待队列的末尾
						WatingQueue = append(WatingQueue, ServingQueue[j])
						// 并且立即服务高风速请求
						ServingQueue[j] = nowReq
						// 跳出查找，同时接下来不进行RR调度
						break
					}
					if j == len(ServingQueue)-1 {
						// 否则，如果服务队列中的优先级都比当前请求高，就将该请求加入等待队列的末尾
						WatingQueue = append(WatingQueue, nowReq)
					}
				}
			}
			// 当处理完请求队列之后，将已处理的元素清空
			RequestQueue = []ScheduleReq{}
		}

		// 处理等待队列和服务队列
		for len(ServingQueue) < MaxServerNum {
			// 如果等待队列为空，就跳出循环了
			if len(WatingQueue) == 0 {
				break
			} else { // 否则就直接往服务队列里调度
				nowReq := WatingQueue[0]
				ServingQueue = append(ServingQueue, nowReq)
				WatingQueue = WatingQueue[1:]
			}
		}

		if SchedulerDB != nil {
			airOrm := repository.AirConditionerOrm{Db: SchedulerDB}
			all, err := airOrm.FindAll()
			if err != nil {
				log.Println(err)
			} else {
				for _, room := range all {
					if room.Power != "on" {
						continue
					}
					fo := repository.FeeOrm{Db: SchedulerDB}
					fee, err2 := fo.QueryByRoomNum(room.RoomNum)
					if err2 != nil {
						log.Println(err)
						fee.RoomNum = room.RoomNum
						fee.Cost = 0
						err := fo.Create(&fee)
						if err != nil {
							log.Println(err)
						}
					} else {
						fee.Cost += float64(model.LevelMap[room.WindLevel]) * model.DefaultFeeRate * 0.005
						err := fo.Update(&fee)
						if err != nil {
							log.Println(err)
						}
						room.TotalFee = fee.Cost
						err = airOrm.Update(room)
						if err != nil {
							log.Println(err)
						}
					}

					if room.Temperature != room.RoomTemperature {
						if math.Abs(room.Temperature - room.RoomTemperature) < 0.5 {
							room.RoomTemperature = room.Temperature
						} else if room.Temperature > room.RoomTemperature {
							room.RoomTemperature += 0.5
						} else {
							room.RoomTemperature -= 0.5
						}
						err := airOrm.Update(room)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
		time.Sleep(time.Second * TimeGap)
	}
}

// 轮询调度
func RoundRobin(){

	for{
		// 如果当前服务队列满
		if len(ServingQueue) == MaxServerNum {
			// 如果等待队列非空，则进行时间片调度
			if len(WatingQueue) != 0 {
				// 等待队列的末尾加上服务队列的第一个元素
				WatingQueue = append(WatingQueue, ServingQueue[0])
				// 将服务队列的第一个元素删除
				ServingQueue = ServingQueue[1:]
				// 服务队列的末尾加上等待队列的第一个元素
				ServingQueue = append(ServingQueue, WatingQueue[0])
				// 将等待队列的第一个元素删除
				WatingQueue = WatingQueue[1:]
			}
		}
		time.Sleep(time.Second*TimeSlice)
	}

}

// 如果请求队列中有关机请求，就移出队列
// 同时打标记，该请求被移除
func checkPowerOff(nowReq ScheduleReq) (flag bool) {
	flag = false
	if nowReq.Power == "off" {
		for i := 0; i < len(ServingQueue); i++ {
			if ServingQueue[i].RoomNum == nowReq.RoomNum {
				ServingQueue = append(ServingQueue[:i], ServingQueue[i+1:]...)
				fmt.Println("checkPowerOff... remove from ServingQueue")
				flag = true
			}
		}

		for i := 0; i < len(WatingQueue); i++ {
			if WatingQueue[i].RoomNum == nowReq.RoomNum {
				WatingQueue = append(WatingQueue[:i], WatingQueue[i+1:]...)
				fmt.Println("checkPowerOff... remove from WatingQueue")
				flag = true
			}
		}
	}
	return flag
}

func checkWindStop(nowReq ScheduleReq) (flag bool) {
	flag = false
	if nowReq.WindLevel == 0 {
		for i := 0; i < len(ServingQueue); i++ {
			if ServingQueue[i].RoomNum == nowReq.RoomNum {
				ServingQueue = append(ServingQueue[:i], ServingQueue[i+1:]...)
				fmt.Println("checkWindStop... remove from ServingQueue")
				flag = true
			}
		}

		for i := 0; i < len(WatingQueue); i++ {
			if WatingQueue[i].RoomNum == nowReq.RoomNum {
				WatingQueue = append(WatingQueue[:i], WatingQueue[i+1:]...)
				fmt.Println("checkWindStop... remove from WatingQueue")
				flag = true
			}
		}
	}
	return flag
}

func checkWindChange(nowReq ScheduleReq) (flag bool) {

	flag = false
	// 如果改变风速，就替换原来的风速等级
	if nowReq.WindFlag == 1 {
		for i := 0; i < len(ServingQueue); i++ {
			if ServingQueue[i].RoomNum == nowReq.RoomNum {
				ServingQueue[i] = nowReq
				fmt.Println("checkWindChange... change from ServingQueue")
				flag = true
			}
		}

		for i := 0; i < len(WatingQueue); i++ {
			if WatingQueue[i].RoomNum == nowReq.RoomNum {
				WatingQueue[i] = nowReq
				fmt.Println("checkWindChange... change from WatingQueue")
				flag = true
			}
		}
	}
	return flag
}

// 开始送风，就修改开始送风数组
func startWind(req ScheduleReq){
	startWindTime := time.Now().Unix()

	val, ok := windTimeMap[req.RoomNum]
	if !ok {
		windTime := WindTime{
			RoomNum: req.RoomNum,
			StartWindTime: []int64{startWindTime},
		}
		windTimeMap[req.RoomNum] = windTime
	} else {
		val.StartWindTime = append(val.StartWindTime, startWindTime)
	}
	
}

// 停止送风，就修改停止送风数组
func stopWind(req ScheduleReq){
	stopWindTime := time.Now().Unix()

	val, ok := windTimeMap[req.RoomNum]
	if !ok {
		windTime := WindTime{
			RoomNum: req.RoomNum,
			StopWindTime: []int64{stopWindTime},
		}
		windTimeMap[req.RoomNum] = windTime
	} else{
		val.StopWindTime = append(val.StopWindTime, stopWindTime)
	}
	
}

func main() {
	// ScheduleReq1 := ScheduleReq{
	// 	WindLevel : 1,
	// }
	// ScheduleReq2 := ScheduleReq{
	// 	WindLevel : 1,
	// }
	// ScheduleReq3 := ScheduleReq{
	// 	WindLevel : 2,
	// }
	// ScheduleReq4 := ScheduleReq{
	// 	WindLevel : 2,
	// }
	// ScheduleReq5 := ScheduleReq{
	// 	WindLevel : 3,
	// }
	// ScheduleReq6 := ScheduleReq{
	// 	WindLevel : 3,
	// }
	// RequestQueue = append(RequestQueue, ScheduleReq1)
	// RequestQueue = append(RequestQueue, ScheduleReq2)
	// RequestQueue = append(RequestQueue, ScheduleReq3)
	// RequestQueue = append(RequestQueue, ScheduleReq4)
	// RequestQueue = append(RequestQueue, ScheduleReq5)
	// RequestQueue = append(RequestQueue, ScheduleReq6)

	Schedule()
}