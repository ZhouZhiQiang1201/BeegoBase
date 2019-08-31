package controllers

import (
	"BeegoBase/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"path"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	userName := this.GetSession("userName")
	fmt.Println(userName)
	//if userName == nil {
	//	this.Redirect("/login", 302)
	//}
	var article []models.Article
	db := orm.NewOrm()
	qs := db.QueryTable("article")
	_, _ = qs.All(&article)
	for _, item := range article{
		_, err := db.LoadRelated(&item, "ArticleType")
		if err != nil {
			fmt.Println(" Read Articles")
			this.TplName = "index.html"
		}
	}


	var articleType []models.ArticleType
	qs1 := db.QueryTable("article_type")
	qs1.All(&articleType)
	this.Data["articles"] = article
	this.Data["articleType"] = articleType
	this.TplName = "index.html"
}

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	name := this.GetString("userName")
	password := this.GetString("password")
	fmt.Println("name: " + name + " password: " + password)
	db := orm.NewOrm()
	user := models.Users{Name: name}
	err := db.Read(&user, "Name")
	if err != nil {
		fmt.Println("Login Error: ", err)
		this.Redirect("/login", 302)
	}
	if user.Password != password {
		fmt.Println("Password Error")
		this.Redirect("/login", 302)
	}
	this.Ctx.SetCookie("userName", name, time.Second * 3600)
	this.SetSession("userName", name)
	this.Redirect("/", 302)
}

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {

	this.TplName = "register.html"
}

func (this *RegisterController) Post() {
	name := this.GetString("username")
	password := this.GetString("password")
	db := orm.NewOrm()
	user := models.Users{}
	user.Name = name
	user.Password = password
	n, err := db.Insert(&user)
	if err != nil {
		fmt.Println("Insert User Error: ", err)
	}
	fmt.Println("Insert User Success: ", n)
	this.Redirect("/login", 302)
}

type TypeController struct {
	beego.Controller
}

func (this *TypeController) Get() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
	}
	var article_type []models.ArticleType
	db := orm.NewOrm()
	qs := db.QueryTable("article_type")
	_, err := qs.All(&article_type)
	if err != nil {
		fmt.Println(" Read Articles")
		this.TplName = "addType.html"
	}
	this.Data["article"] = article_type
	this.TplName = "addType.html"
}

func (this *TypeController) Post() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
	}
	name := this.GetString("typeName")
	if name == "" {
		fmt.Println(" add Type name is Null")
		this.Redirect("/addType", 302)
	}
	db := orm.NewOrm()
	article_type := models.ArticleType{}
	article_type.Name = name
	_, err := db.Insert(&article_type)
	if err != nil {
		fmt.Println("Insert Article Type Error:", err)
	}
	this.Redirect("/addType", 302)
}

func (this *TypeController) DeleteType() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
	}
	typeId := this.GetString("Id")
	articleTypeId, _ := strconv.Atoi(typeId)
	db := orm.NewOrm()
	ArticleType := models.ArticleType{Id: articleTypeId}
	_, err := db.Delete(&ArticleType)
	if err != nil {
		fmt.Println("Read Error")
	}
	this.Redirect("/addType", 302)

}

type NewController struct {
	beego.Controller
}

func (this *NewController) Get() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
	}
	var ArticleType []models.ArticleType
	db := orm.NewOrm()
	qs := db.QueryTable("article_type")
	_, err := qs.All(&ArticleType)
	if err != nil {
		fmt.Println("Read Article Type Error!")
		this.TplName = "add.html"
	}
	this.Data["ArticleType"] = ArticleType
	this.TplName = "add.html"
}

func (this *NewController) Post() {
	title := this.GetString("articleName")
	content := this.GetString("content")
	typeId := this.GetString("select")
	if content == "" || title == "" {
		this.Data["errmsg"] = "参数缺失"
		return
	}
	f, head, err := this.GetFile("uploadname")
	defer f.Close()
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}

	//文件大小
	if head.Size >= 1024*1024*50 {
		this.Data["errmsg"] = "file size error"
		this.TplName = "add.html"
	}
	// 文件后缀
	ext := path.Ext(head.Filename)
	fmt.Println("ext: ", ext)
	if ext != ".png" && ext != "jpg" && ext != "jpeg" {
		this.Data["errmsg"] = "格式错误"
		this.TplName = "add.html"
	}
	fileName := time.Now().Format("20060102-150405") + ext
	//fileName := head.Filename
	err = this.SaveToFile("uploadname", "./static/img/"+fileName)
	if err != nil {
		fmt.Println("Save File Error: ", err)
		this.Data["errmsg"] = "Save File Error"
		return
	}

	db := orm.NewOrm()
	type_id, _ := strconv.Atoi(typeId)
	articleType := models.ArticleType{Id: type_id}
	err = db.Read(&articleType)
	if err != nil {
		fmt.Println("ll")
		return
	}
	article := models.Article{}
	article.Title = title
	article.Content = content
	article.Image = fileName
	article.ArticleType = &articleType
	_, err = db.Insert(&article)
	if err != nil {
		fmt.Println("Insert Article Error: ", err)
		this.Redirect("/", 302)
	}
	this.Redirect("/", 302)
}

func (this *NewController) DeleteArticle() {
	articleId, _ := this.GetInt("Id")
	db := orm.NewOrm()
	Article := models.Article{Id:articleId}
	db.Delete(&Article)
	this.Redirect("/", 302)


}

func (this *NewController) GetArticle() {
	articleId, _ := this.GetInt("Id")
	db := orm.NewOrm()
	Article := models.Article{Id:articleId}
	db.Read(&Article)
	this.Data["Article"] = Article
	this.TplName = "update.html"
}

func (this *NewController) UpdateArticle() {
	articleId, _ := this.GetInt("Id")
	content := this.GetString("content")
	title := this.GetString("articleName")
	fmt.Println(articleId, content)
	db := orm.NewOrm()
	Article := models.Article{Id:articleId}
	_ = db.Read(&Article)
	Article.Title = title
	Article.Content = content
	_, err := db.Update(&Article)
	if err != nil {
		fmt.Println("Update Article Error", err)
	}

	this.Redirect("/", 302)
}

type ContentController struct {
	beego.Controller
}

func (this *ContentController) Get() {
	ArticleId := this.Ctx.Input.Param(":Id")
	fmt.Println("articleId", ArticleId)

	articleId, err := strconv.Atoi(ArticleId)
	if err != nil {
		fmt.Println("ga")
		this.Data["errmsg"] = "转换失败"
		return
	}
	fmt.Printf("article type: %T", ArticleId)
	Article := models.Article{Id: articleId}
	//article.Id = articleId
	db := orm.NewOrm()
	err = db.Read(&Article)
	if err != nil {
		fmt.Println("read error")
		this.Data["errmsg"] = "Read Error"
		return
	}

	this.Data["Title"] = Article.Title
	this.Data["CreateTime"] = Article.CreateTime
	this.Data["Content"] = Article.Content
	this.Data["Count"] = Article.ReadCount
	this.Data["image"] = "/static/img/" + Article.Image
	Article.ReadCount += 1
	db.Update(&Article)
	this.TplName = "content.html"
}
