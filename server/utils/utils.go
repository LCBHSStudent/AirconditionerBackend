package utils

import (
	"AirConditioner/server/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	Conn   net.Conn
	Buffer [4096]byte
}

func (this *Transfer) ReadPkg() (msg message.Message, err error) {
	buf := this.Buffer
	conn := this.Conn
	// 从conn中读取字节，放到buf中
	_, err = conn.Read(buf[:4096])
	if err != nil {
		return
	}
	// buf反序列化
	err = json.Unmarshal(buf[:4096], &msg)
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4096])
	if int(pkgLen) != msg.Length {
		fmt.Println("pkg Length wrong...")
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	conn := this.Conn
	pkgLen := uint32(len(data))
	n, err := conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail, err=", err)
		return
	}
	return
}
