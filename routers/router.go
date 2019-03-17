package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/paylm/myweb/models/user"
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
	fmt.Printf("loadUserInfo u:%v\n", u)
	if u == nil {
		return user.UserData{Username: "anon", Job: "vister", Email: ""}
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
	pjs := user.GetProjects()
	u := loadUserInfo(c)
	fmt.Printf("userinfo:%v\n", u)
	c.HTML(200, "index.html", gin.H{
		"title":  "index",
		"online": online,
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
	session.Set("user", u)
	session.Save()
	//c.String(http.StatusOK, "Login successful")
	c.Redirect(http.StatusMovedPermanently, "/")
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
			c.Redirect(http.StatusMovedPermanently, "/")
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
	c.HTML(200, "tables.html", gin.H{
		"title": "ppl table",
		"msg":   "",
		"users": users,
	})
}
func charts(c *gin.Context) {
	c.HTML(200, "charts.html", gin.H{
		"title": "ppl charts",
		"msg":   "",
	})
}

func logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":   200,
		"result": "logout ok",
	})
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

func foo(c *gin.Context) {
	c.JSON(200, "foo")
}

func InitRouter() *gin.Engine {
	r := gin.New()

	//init session
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})

	r.Use(sessions.Sessions("mysession", store))

	r.Use(gin.Recovery())
	r.Static("/static", "./template")
	r.LoadHTMLGlob("template/*.html")
	gin.SetMode(setting.ServerSetting.RunMode)

	r.GET("/", index)
	r.GET("/logout", logout)
	r.POST("/login", login)
	r.GET("/login", loginPage)
	r.GET("/tables", tables)
	r.GET("/charts", charts)
	r.Any("/register", register)
	bl := r.Group("/blog")
	{
		bl.GET("/list", blist)
		bl.POST("/edit", bedit)
	}

	return r
}
