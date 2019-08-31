package main

import (
	_ "BeegoBase/models"
	_ "BeegoBase/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
