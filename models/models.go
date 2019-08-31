package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"time"
)

type Users struct {
	Id       int
	Name     string
	Password string
}

type ArticleType struct {
	Id int
	Name string
	Article []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id int
	Title string
	Content string
	Image string
	CreateTime time.Time `orm:"auto_now"`
	ReadCount int `orm:"default(0)"`

	ArticleType *ArticleType `orm:"rel(fk)"`
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", "root:mysql@tcp(192.168.101.14:3306)/beego?charset=utf8")
	if err != nil {
		fmt.Println("db error: ", err)
	}
	orm.RegisterModel(new(Users), new(ArticleType), new(Article))
	err  = orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println("Create Users Tables Error: ", err)
	}


}
