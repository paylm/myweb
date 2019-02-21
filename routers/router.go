package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paylm/myweb/models/user"
	"github.com/paylm/myweb/pkg/gredis"
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
	c.HTML(200, "index.html", gin.H{
		"title": "index",
	})
}

func login(c *gin.Context) {
	var u user.UserData
	err := c.Bind(&u)
	if c.Bind(&u) != nil {
		fmt.Printf("login bind error :%v\n", err)
	}
	err = gredis.Set(fmt.Sprintf("login:%s", u.Username), 1, 3600)
	if err != nil {
		fmt.Printf("login redis err:%v\n", err)
	}
	resLogin := u.Verlogin()
	if resLogin != nil {
		c.HTML(200, "login.html", gin.H{
			"title": "注册",
			"msg":   fmt.Sprintf("%s", resLogin),
		})

		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "123",
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)
	c.String(http.StatusOK, "Login successful")
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
	c.HTML(200, "tables.html", gin.H{
		"title": "ppl table",
		"msg":   "",
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
