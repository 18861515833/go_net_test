package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"../gen/pb"
	"github.com/golang/protobuf/proto"
)

//channel 只声明不make  使用是会出问题的
//var quitSemaphore chan bool
var (
	quitSemaphore = make(chan bool)

	name string
)

func main() {
	//先创建addr
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	//连接
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	fmt.Println("connected!")
	fmt.Print("输入昵称：")
	fmt.Scanln(&name)

	//创建两个线程，一个收，一个发送
	go sendThread(conn)
	go onMessageRecived(conn)

	//阻塞主线程
	res := <-quitSemaphore
	if res {
		fmt.Println("客户端异常，程序退出")
	} else {
		fmt.Println("服务器异常，程序退出")
	}
}
func sendThread(conn *net.TCPConn) {
	//var msg message = message{Name: name}
	var msg =&pb.Message{Name:name}
	input := bufio.NewScanner(os.Stdin)
	headbuf := make([]byte, 4)
	for input.Scan() {
		//读取消息内容
		msg.Content = input.Text()
		//填充bodybuf
		//bodybuf, err := json.Marshal(msg)
		bodybuf,err:=proto.Marshal(msg)
		if err != nil {
			//fmt.Println("json序列化失败")
			fmt.Println("protobuf序列化失败")
			return
		}
		//求出body的len
		msglen := len(bodybuf)
		//填充headuf
		binary.BigEndian.PutUint32(headbuf, uint32(msglen))

		if err != nil {
			return
		}
		conn.Write(headbuf)
		conn.Write(bodybuf)
		//fmt.Println("head:",headbuf,"headlen:",4)
		//fmt.Println("body:",bodybuf,"bodylen:",msglen)
	}
	quitSemaphore <- true
}

func onMessageRecived(conn *net.TCPConn) {
	headbuf := make([]byte, 4)
	datalen := uint32(0)
	for {
		_, err := io.ReadFull(conn, headbuf)
		if err != nil {
			break
		}
		datalen = binary.BigEndian.Uint32(headbuf)
		data := make([]byte, datalen)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			break
		}
		//var msg message
		//json.Unmarshal(data, &msg)
		var msg pb.Message
		proto.Unmarshal(data,&msg)
		fmt.Println(msg.Name, ":", msg.Content)
	}
	quitSemaphore <- false
}
