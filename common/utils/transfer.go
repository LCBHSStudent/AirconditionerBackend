package utils

import (

	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"net"
)

// 将这些方法关联到结构体中
type Transfer struct {
	// 分析应该有哪些字段
	Conn net.Conn
	Buf  [4096]byte // 传输时使用的缓冲
}

// 将读取数据包，直接封装成一个函数readPkg(),返回message，err
func (this *Transfer) ReadPkg() (msg message.Message, err error) {

	// buf := make([]byte, 1024*4)
	buf := this.Buf
	conn := this.Conn
	// fmt.Println("读取客户端发送的数据...")
	// conn.Read只有在conn没有被关闭的情况下才会阻塞
	// 如果客户端关闭了conn则不阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//fmt.Println("读取到的buf=",buf[:4])

	// 根据buf[:4]转换成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])

	// 根据pkgLen读取消息内容，从conn中读取字节，放到buf中
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		return
	}

	// 把buf反序列化为message.Message
	// 技术就是一层窗户纸
	err = json.Unmarshal(buf[:pkgLen], &msg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {

	// 先发送长度给对方
	// 先获取到data的长度，然后转成一个表示长度的byte切片
	buf := this.Buf
	conn := this.Conn
	var pkgLen uint32
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	// 发送长度
	n, err := conn.Write(buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}

	// 发送data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	return
}