package main

import(
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
	"xorm.io/core"
)

var engine *xorm.Engine

type User struct {
	Id int64
	Name string `xorm:"varchar(25) not null unique 'usr_name'"`
	CreateAt time.Time `xorm:"created"`
	GroupId int64 `xorm:index`
}

type Group struct {
	Id int64
	Name string
}

func main() {
	var err error
	engine,err=xorm.NewEngine("mysql","root:123456@(47.102.121.34:3306)/test?charset=utf8")
	if err!=nil{
		fmt.Println(err.Error())
	}
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	engine.SetMapper(core.SameMapper{})
	engine.Sync()

	engine.CreateTables(User{})

	engine.Close()
}
