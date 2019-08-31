package routers

import (
	"BeegoBase/controllers"
	_ "BeegoBase/models"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/add", &controllers.NewController{})
	beego.Router("/delete", &controllers.NewController{}, "get:DeleteArticle")
	beego.Router("/update", &controllers.NewController{}, "get:GetArticle;post:UpdateArticle")
	beego.Router("/addType", &controllers.TypeController{})
	beego.Router("/deleteType", &controllers.TypeController{}, "get:DeleteType")
	beego.Router("/content/:Id", &controllers.ContentController{})
}
