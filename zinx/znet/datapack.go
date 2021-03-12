package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

// DataPack 封包拆包的具体模块
type DataPack struct {
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节）+ ID uint32（4字节）
	return 8
}

// Pack 封包
func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	// 创建一个存放字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将DataLen写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetDataLen()); err != nil {
		return nil, err
	}

	// 将ID写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetMsgID()); err != nil {
		return nil, err
	}

	// 将Data数据写入dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包，只需要将head信息读取出来，再根据head信息中消息内容的长度进行一次读
func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	// 创建一个读取二进制数据的ioReader
	dataBuff := bytes.NewReader(data)

	// 只读取head信息，得到DataLen和ID
	message := &Message{}

	// 读DataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &message.DataLen); err != nil {
		return nil, err
	}
	// 读ID
	if err := binary.Read(dataBuff, binary.LittleEndian, &message.ID); err != nil {
		return nil, err
	}

	// 判断DataLen是否已经超出了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && message.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large message data received")
	}
	return message, nil
}

// NewDataPack 初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}
