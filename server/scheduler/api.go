package scheduler

func AddScheduleReq(req ScheduleReq) {
	RequestQueue = append(RequestQueue, req)
}

func GetServingQueue() (roomQueue []int) {
	for i:=0 ; i < len(ServingQueue); i++ {
		roomQueue = append(roomQueue, ServingQueue[i].RoomNum)
	}
	return roomQueue
}