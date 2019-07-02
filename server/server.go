package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	//"time"
)

var (
	client =make([] *net.TCPConn,0,100)
)

func main1(){
	data:=[5]int{1,2,3,4,5}
	s1:=data[:]
	s2:=data[1:]
	data[4]=0
	fmt.Println(data)
	fmt.Println(s1)
	fmt.Println(s2)
}
func init(){
	fmt.Println("init")
	fmt.Print()
}

func main(){
	var tcpAddr *net.TCPAddr

	//resolve 解析
	tcpAddr,_=net.ResolveTCPAddr("tcp","127.0.0.1:9999")

	tcpListener,_:=net.ListenTCP("tcp",tcpAddr)

	//推迟执行，防止资源泄露
	defer tcpListener.Close()

	//while(true) for(;;)
	for {
		tcpConn,err:=tcpListener.AcceptTCP()
		//如果发生错误，继续监听，不用管
		if err!= nil{
			continue
		}
		client=append(client,tcpConn)
		fmt.Println("client:",len(client))
		//客户端正常连接
		fmt.Println("a client connected:"+tcpConn.RemoteAddr().String())
		//为客户端创建单独线程
		go tcpPipe(len(client)-1)
	}
}

func tcpPipe(index int){
	conn:=client[index]
	ipStr:=conn.RemoteAddr().String()
	
	//推迟客户端的退出，防止资源泄露
	defer func(){
		fmt.Println("disconnected :"+ipStr)
		conn.Close()
		client=append(client[:index],client[index+1:]...)
	}()

	//reader := bufio.NewReader(conn)
	headbuf:=make([]byte,4)
	datalen:=uint32(0)
	for {
		_,err:=io.ReadFull(conn,headbuf)
		if err !=nil {
			break;
		}
		datalen=binary.BigEndian.Uint32(headbuf)
		data:=make([]byte,datalen)
		_,err=io.ReadFull(conn,data)
		if err != nil {
			break
		}
		fmt.Println(data)
		//给client中的每个客户端都发送消息
		for i:=0;i<len(client);i++ {
			client[i].Write(headbuf)
			client[i].Write(data)
		}
	}
	/*
	for{
		message,err:=reader.ReadString('\n')
		//如果出错，直接return
		if err !=nil{
			return
		}
		//打印recv的消息
		fmt.Println(string(message))
		//重新包装recv的消息
		//msg:=time.Now().String()+"\n"
		mes:="聊天消息："+string(message)
		b:=[]byte(mes)
		//发送包装完成的消息
		//conn.Write(b)
		
		//给client中的每个客户端都发送消息
		for i:=0;i<len(client);i++ {
			client[i].Write(b)
		}
	}
	*/
}