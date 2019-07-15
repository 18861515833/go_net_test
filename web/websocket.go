package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func echo(writer http.ResponseWriter, request *http.Request) {

	//升级websocket
	conn,err :=(&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(writer,request,nil)
	if err!=nil{
		fmt.Println(err.Error())
		//return
	}
	for {
		_,data,err := conn.ReadMessage()
		if err!=nil{
			fmt.Println()
		}
		//这里没有使用json，所以不用解析，直接把byte数组转化成string即可
		msg:=string(data)
		fmt.Println(msg)
	}
}

func main() {

	http.HandleFunc("/", echo)
	http.ListenAndServe(":1234",nil)

}
