package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"net"
)

type AirProcessor struct {
	Conn net.Conn
	Orm  *repository.AirConditionerOrm
}

// 查询空调状态，按房间号查询
func (ap *AirProcessor) FindByRoom(msg *message.Message) (err error) {
	var findByRoom message.AirConditionerFindByRoom
	err = json.Unmarshal([]byte(msg.Data), &findByRoom)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var findByRoomRes message.AirConditionerFindByRoomRes

	roomNum := findByRoom.RoomNum
	airConditioner, err := ap.Orm.FindByRoom(roomNum)
	if err != nil {
		fmt.Println("ap.Orm.FindByRoom err=", err)
		return err
	}

	findByRoomRes.Code = 200
	findByRoomRes.Msg = "根据房间号查询空调成功！"
	findByRoomRes.AirConditioner = airConditioner

	data, err := json.Marshal(findByRoomRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindByRoomRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// FindAll 查询所有空调状态
func (ap *AirProcessor) FindAll(msg *message.Message) (err error) {

	var resMsg message.Message
	var findAllRes message.AirConditionerFindAllRes

	airConditioners, err := ap.Orm.FindAll()
	if err != nil {
		fmt.Println("ap.Orm.FindAll() err=", err)
		return err
	}

	findAllRes.Code = 200
	findAllRes.Msg = "查询所有空调状态成功！"
	findAllRes.AirConditioners = airConditioners

	data, err := json.Marshal(findAllRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindAllRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// Create 新增一条空调状态记录
func (ap *AirProcessor) Create(msg *message.Message) (err error) {
	var createMsg message.AirConditionerCreate
	err = json.Unmarshal([]byte(msg.Data), &createMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var CreateRes message.NormalRes

	airConditioner := createMsg.AirConditioner

	existAir, _ := ap.Orm.FindByRoom(airConditioner.RoomNum)
	// 如果已存在，则返回重新创建提示
	if existAir.RoomNum != 0 {
		CreateRes.Code = 501
		CreateRes.Msg = "空调Room已存在，请重新创建！"
	} else {
		err = ap.Orm.Create(&airConditioner)
		if err != nil {
			fmt.Println("ap.Orm.Create(&airConditioner) err=", err)
			return
		}
		CreateRes.Code = 200
		CreateRes.Msg = "创建成功！"
	}

	data, err := json.Marshal(CreateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) Update(msg *message.Message) (err error) {
	var updateMsg message.AirConditionerUpdate
	err = json.Unmarshal([]byte(msg.Data), &updateMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var updateRes message.NormalRes

	airConditioner := updateMsg.AirConditioner
	err = ap.Orm.Update(airConditioner)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	updateRes.Code = 200
	updateRes.Msg = "更新成功！"

	data, err := json.Marshal(updateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// AirConditionerOn
func (ap *AirProcessor) PowerOn(msg *message.Message) (err error) {
	var powerOn message.AirConditionerOn
	err = json.Unmarshal([]byte(msg.Data), &powerOn)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(powerOn.RoomNum)
	if err != nil {
		fmt.Println("ap.Orm.FindByRoom(powerOn.RoomNum) err = ", err)
		return err
	}
	// 这里需要对 OpenTime 进行处理
	if air.OpenTime == "" { // 如果是第一次开机，就初始化 OpenTime
		openTimeList := []int64{powerOn.OpenTime}
		res, err := json.Marshal(openTimeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.OpenTime = string(res)
	} else { // 否则，就append OpenTime 列表
		var timeList []int64
		json.Unmarshal([]byte(air.OpenTime), &timeList)
		timeList = append(timeList, powerOn.OpenTime)
		res, err := json.Marshal(timeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.OpenTime = string(res)
	}
	air.Power = model.PowerOn
	air.Mode = powerOn.Mode
	air.WindLevel = powerOn.WindLevel
	air.Temperature = powerOn.Temperature

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "开机成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// AirConditionerOff
func (ap *AirProcessor) PowerOff(msg *message.Message) (err error) {
	var powerOff message.AirConditionerOff
	err = json.Unmarshal([]byte(msg.Data), &powerOff)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(powerOff.RoomNum)
	if err != nil {
		return err
	}
	// 这里需要对 CloseTime 进行处理
	if air.CloseTime == "" { // 如果是第一次开机，就初始化 OpenTime
		timeList := []int64{powerOff.CloseTime}
		res, err := json.Marshal(timeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.CloseTime = string(res)
	} else { // 否则，就append OpenTime 列表
		var timeList []int64
		json.Unmarshal([]byte(air.CloseTime), &timeList)
		timeList = append(timeList, powerOff.CloseTime)
		res, err := json.Marshal(timeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.CloseTime = string(res)
	}
	air.Power = model.PowerOff
	fmt.Println("air.Power:", air.Power)
	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "关机成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// SetParam
func (ap *AirProcessor) SetParam(msg *message.Message) (err error) {
	var setParam message.AirConditionerSetParam
	err = json.Unmarshal([]byte(msg.Data), &setParam)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(setParam.RoomNum)
	if err != nil {
		return err
	}

	// 设置相应参数
	air.Mode = setParam.Mode
	air.WindLevel = setParam.WindLevel
	air.Temperature = setParam.Temperature
	// 调整次数加一
	air.SetParamNum += 1

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "调整空调参数成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) SetRoomData(msg *message.Message) (err error) {
	var setRoomData message.SetRoomData
	err = json.Unmarshal([]byte(msg.Data), &setRoomData)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(setRoomData.RoomNum)
	if err != nil {
		return err
	}

	// 设置相应参数
	air.TotalPower = setRoomData.TotalPower
	air.RoomTemperature = setRoomData.RoomTemperature

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "设置空调耗电量和房间温度成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// StopWind
func (ap *AirProcessor) StopWind(msg *message.Message) (err error) {
	var stopWind message.AirConditionerStopWind
	err = json.Unmarshal([]byte(msg.Data), &stopWind)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(stopWind.RoomNum)
	if err != nil {
		return err
	}

	// 这里需要对 stopWind 进行处理
	if air.StopWind == "" { // 如果是第一次开机，就初始化 OpenTime
		timeList := []int64{stopWind.StopWindTime}
		res, err := json.Marshal(timeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.StopWind = string(res)
	} else { // 否则，就append OpenTime 列表
		var timeList []int64
		json.Unmarshal([]byte(air.StopWind), &timeList)
		timeList = append(timeList, stopWind.StopWindTime)
		res, err := json.Marshal(timeList)
		if err != nil {
			fmt.Println("json.Marshal(timeList) err =",err)
		}
		air.StopWind = string(res)
	}

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "空调停止送风成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// GetReport
func (ap *AirProcessor) GetReport(msg *message.Message) (err error) {

	var resMsg message.Message
	var getReportRes message.GetReportRes

	airs, err := ap.Orm.FindAll()
	if err != nil {
		fmt.Println("ap.Orm.FindAll() err=", err)
		return err
	}

	var (
		openTimeList []int64
		closeTimeList []int64
	)


	for _, air := range airs {
		json.Unmarshal([]byte(air.OpenTime), &openTimeList)
		json.Unmarshal([]byte(air.CloseTime), &closeTimeList)
		var report model.Report
		report.RoomNum = air.RoomNum
		report.TotalPower = air.TotalPower
		report.TotalFee = air.TotalPower * 1
		report.CloseNum = len(closeTimeList)
		report.SetParamNum = air.SetParamNum

		for i := 0; i < len(closeTimeList); i++ {
			// 计算空调的总开机时长：用关机数组的值逐个减去开机数组的值
			report.UsedTime += int(air.CloseTime[i] - air.OpenTime[i])
		}

		// 将空调报表逐个添加到 getReportRes.Reports 中
		getReportRes.Reports = append(getReportRes.Reports, report)
	}

	getReportRes.Code = 200
	getReportRes.Msg = "获取所有空调报表成功！"

	data, err := json.Marshal(getReportRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeGetReportRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// GetDetailList
func (ap *AirProcessor) GetDetailList(msg *message.Message) (err error) {

	var getDetailList message.GetDetailList
	err = json.Unmarshal([]byte(msg.Data), &getDetailList)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}
	
	air, err := ap.Orm.FindByRoom(getDetailList.RoomNum)
	if err != nil {
		fmt.Println("Orm.FindByRoom err=", err)
		return err
	}

	var resMsg message.Message
	var getDetailListRes message.GetDetailListRes

	var (
		startWindList []int64
		stopWindList []int64
	)

	json.Unmarshal([]byte(air.StartWind), &startWindList)
	json.Unmarshal([]byte(air.StopWind), &stopWindList)
	var detail model.Detail
	detail.RoomNum = air.RoomNum
	detail.StartWindList = startWindList
	detail.StoptWindList = stopWindList
	detail.WindLevel = air.WindLevel
	detail.TotalPower = air.TotalPower
	detail.FeeRate = model.DefaultFeeRate
	detail.TotalFee = air.TotalPower * model.DefaultFeeRate

	for i := 0; i < len(startWindList); i++ {
		// 计算空调的总送风时长：用停止送风数组的值逐个开始送风数组的值
		detail.TotalWindTime += int64(stopWindList[i] - startWindList[i])
	}

	getDetailListRes.Code = 200
	getDetailListRes.Msg = "获取详单成功！"

	data, err := json.Marshal(getDetailListRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeGetDetailListRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}
