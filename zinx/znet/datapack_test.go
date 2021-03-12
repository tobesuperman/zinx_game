package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只负责DataPack封包拆包的单元测试
func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Server listen err:", err)
		return
	}
	// 创建一个goroutine承载负责从客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Server accept err:", err)
				continue
			}
			go func(conn net.Conn) {
				// 处理客户端的请求，拆包的过程
				dp := NewDataPack()
				for {
					// 第一次读head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("Read msg head err:", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("Server unpack err:", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						// 第二次读消息内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("Server unpack err:", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println("Receive ID =", msg.ID, "DataLen =", msg.DataLen, "Data =", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Client dial error:", err)
		return
	}

	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg包一同发送
	// 封装第一个msg包
	data1 := []byte("zinx")
	msg1 := &Message{
		ID:      0,
		DataLen: uint32(len(data1)),
		Data:    data1,
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Client pack msg1 error:", err)
		return
	}

	data2 := []byte("hello, go")
	msg2 := &Message{
		ID:      1,
		DataLen: uint32(len(data2)),
		Data:    data2,
	}
	sendData2, err := dp.Pack(msg2)

	data3 := []byte("hello, zinx")
	msg3 := &Message{
		ID:      1,
		DataLen: uint32(len(data3)),
		Data:    data3,
	}
	sendData3, err := dp.Pack(msg3)
	if err != nil {
		fmt.Println("Client pack msg2 error:", err)
		return
	}

	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	sendData1 = append(sendData1, sendData3...)
	// 一次性发给服务器
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
