package routers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/paylm/myweb/models/user"
	"github.com/paylm/myweb/pkg/middleware"
	"github.com/paylm/myweb/pkg/setting"
)

type b struct {
	b1 string
}

type blog struct {
	title   string
	content string
	pubtime string
}

func loadUserInfo(c *gin.Context) interface{} {

	session := sessions.Default(c)
	u := session.Get("user")
	fmt.Printf("loadUserInfo uid:%v\n", session.Get("userId"))
	fmt.Printf("loadUserInfo u:%v\n", u)
	if u == nil {
		return user.UserData{Username: "anon", Job: "vister", Img: "avatar-1.jpg", Email: ""}
	}
	return u
}

func getb(c *gin.Context) {
	var b1 b
	err := c.Bind(&b1)
	if err != nil {
		c.String(200, `getb fail`)
		return
	}
	c.JSON(200, "wellcome to ")
}

func index(c *gin.Context) {
	//c.String(200, `index`)
	online := user.OnlineCount()
	active := user.RecentActive()
	pjs := user.GetProjects()
	u := loadUserInfo(c)
	fmt.Printf("userinfo:%v\n,pjs:%v\n", u, pjs)
	log.Printf("====== index  ======\n")
	c.HTML(200, "index.html", gin.H{
		"title":  "index",
		"online": online,
		"active": active,
		"pjs":    pjs,
		"user":   u,
	})
}

func login(c *gin.Context) {
	var u user.UserData
	err := c.Bind(&u)
	if c.Bind(&u) != nil {
		fmt.Printf("login bind error :%v\n", err)
	}
	u, resLogin := u.Verlogin()
	if resLogin != nil {
		c.HTML(200, "login.html", gin.H{
			"title": "登录",
			"msg":   fmt.Sprintf("%s", resLogin),
		})

		return
	}
	fmt.Printf("%v login ok\n", u)
	//set session
	session := sessions.Default(c)
	session.Set("userId", u.ID)
	session.Set("user", u)
	session.Save()
	//c.String(http.StatusOK, "Login successful")
	c.Redirect(http.StatusFound, "/")
}

func register(c *gin.Context) {
	Username := c.PostForm("username")
	if Username == "" {
		c.HTML(200, "register.html", gin.H{
			"title": "注册",
			"msg":   "信息不能为空",
		})
	}

	var u user.UserData
	if c.Bind(&u) == nil {
		if err := u.Reg(); err == nil {
			c.Redirect(http.StatusFound, "/")
		} else {
			c.HTML(200, "register.html", gin.H{
				"title": "注册",
				"msg":   fmt.Sprintf("%v", err),
			})
		}
	} else {
		fmt.Println("load fail")
	}

}

func loginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"title": "ppl登录",
		"msg":   "",
	})
}

func tables(c *gin.Context) {
	users := user.GetAllUser(100)
	u := loadUserInfo(c)
	c.HTML(200, "tables.html", gin.H{
		"title": "ppl table",
		"msg":   "",
		"users": users,
		"user":  u,
	})
}
func charts(c *gin.Context) {
	u := loadUserInfo(c)
	c.HTML(200, "charts.html", gin.H{
		"title": "ppl charts",
		"msg":   "",
		"user":  u,
	})
}

func forms(c *gin.Context) {
	u := loadUserInfo(c)
	c.HTML(200, "forms.html", gin.H{
		"title": "ppl forms",
		"msg":   "",
		"user":  u,
	})
}

func bookProject(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("book id %d\n", user.BookPj()))
}

func bookUnlockProject(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("book id %d\n", user.BookUoLockPj()))
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("userId")
	//fmt.Printf("logout uid:%v\n", uid)
	if uid == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user.Logout(uid.(int))
	session.Delete("userId")
	session.Delete("user")
	session.Save()
	//	c.JSON(200, gin.H{
	//		"code":   200,
	//		"result": "logout ok",
	//	})
	c.Redirect(http.StatusFound, "/login")
}

func stats(c *gin.Context) {
	aType := c.Param("type")
	//fmt.Printf("stat param :%s\n", aType)
	data := user.FindByType(aType)
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":  http.StatusOK,
			"error": "null",
			"data":  data,
		},
	)
}

func params(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"page": page,
		"name": name,
	},
	)
}

func blist(c *gin.Context) {
	data := make(map[string]blog)
	data["love"] = blog{title: "my love", content: "sssssssss"}
	data["like"] = blog{title: "my like", content: "liek"}
	data["fuck"] = blog{title: "my fuck", content: "liek"}
	c.JSON(200, gin.H{
		"list": data,
	})
}

func bedit(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":   200,
		"result": "logout ok",
	})
}

func uploads(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func convert(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, fmt.Sprintf("./uploads/%s", file.Filename))
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func foo(c *gin.Context) {
	c.JSON(200, "foo")
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}

func InitRouter() *gin.Engine {
	gin.SetMode(setting.ServerSetting.RunMode)
	r := gin.New()
	//init session
	//store := sessions.NewCookieStore([]byte("secret"))
	//http://127.0.0.1:6060/pkg/github.com/gorilla/sessions/
	//As it's not possible to pass a raw type as a parameter to a function, gob.Register() relies on us passing it a value of the desired type
	gob.Register(&user.UserData{})
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	//gob.Register(&{})
	store, _ := sessions.NewRedisStore(10, "tcp", setting.RedisSetting.Host, setting.RedisSetting.Password, []byte("secret"))
	//store.Options(sessions.Options{
	//	MaxAge: int(30 * time.Minute), //30min
	//	Path:   "/",
	//})

	r.Use(sessions.Sessions("mysession", store))

	r.Use(gin.Recovery())
	r.Use(middleware.Logger()) //使用中间件
	r.Static("/static", "./template")
	r.LoadHTMLGlob("template/*.html")

	r.GET("/", index)
	r.GET("/ping", ping)
	r.GET("/logout", logout)
	r.POST("/login", login)
	r.GET("/login", loginPage)
	r.GET("/tables", tables)
	r.GET("/charts", charts)
	r.GET("/forms", forms)
	r.Any("/register", register)
	r.GET("/bookpj", bookProject)
	r.GET("/bookpj2", bookUnlockProject)
	r.GET("/stat/:type", stats)
	r.Any("/params", params)
	r.POST("uploads", uploads)
	bl := r.Group("/blog")
	{
		bl.GET("/list", blist)
		bl.POST("/edit", bedit)
	}

	return r
}
