package scheduler

type ScheduleReq struct {
	RoomNum int
	Power string		// 空调电源："on"/"off"
	WindLevel int		// 空调风速，对应优先级："stop":0,""low":1, "mid":2, "high":3
	ArivingTime int64   // 请求到达的时间
	WindFlag	int 	// 表示是否改变风速，0不改变，1改变，默认为0不改变
}

