package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myzinx/utils"
	"myzinx/ziface"
)

// 封包，拆包的具体模块

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4 字节) + ID uint32(4 字节)
	return 8
}

// 封包方法
// /datalen/msgID/data/
func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法 将包的Head信息读出来 之后更具head信息里的data的长度，再进行一次读
func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到datalen和MsgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen是否超出我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv")
	}

	return msg, nil
}
